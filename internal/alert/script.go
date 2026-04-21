package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
)

// ScriptHandler runs an external script for each alert change.
type ScriptHandler struct {
	cfg config.ScriptConfig
}

// NewScriptHandler creates a ScriptHandler from the given config.
func NewScriptHandler(cfg config.ScriptConfig) *ScriptHandler {
	return &ScriptHandler{cfg: cfg}
}

// Handle executes the configured script once per change, passing a JSON
// payload via stdin.
func (h *ScriptHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		if err := h.run(c); err != nil {
			return err
		}
	}
	return nil
}

// Drain is a no-op for ScriptHandler.
func (h *ScriptHandler) Drain() error { return nil }

func (h *ScriptHandler) run(c monitor.Change) error {
	payload, err := buildScriptPayload(c)
	if err != nil {
		return fmt.Errorf("script: marshal payload: %w", err)
	}

	timeout := h.cfg.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	args := append([]string{}, h.cfg.Args...)
	cmd := exec.CommandContext(ctx, h.cfg.Path, args...)
	cmd.Env = append(cmd.Environ(), h.cfg.EnvVars...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("script: stdin pipe: %w", err)
	}
	go func() {
		defer stdin.Close()
		stdin.Write(payload) //nolint:errcheck
	}()

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("script %s: %w: %s", h.cfg.Path, err, string(out))
	}
	return nil
}

func buildScriptPayload(c monitor.Change) ([]byte, error) {
	type payload struct {
		Kind    string `json:"kind"`
		Proto   string `json:"proto"`
		Addr    string `json:"addr"`
		Port    int    `json:"port"`
		Process string `json:"process,omitempty"`
		PID     int    `json:"pid,omitempty"`
	}
	p := payload{
		Kind:    string(c.Kind),
		Proto:   c.Binding.Proto,
		Addr:    c.Binding.Addr,
		Port:    c.Binding.Port,
		Process: c.Binding.Process,
		PID:     c.Binding.PID,
	}
	return json.Marshal(p)
}
