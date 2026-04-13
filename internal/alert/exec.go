package alert

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/user/portwatch/internal/monitor"
)

// ExecHandler runs an external command for each batch of changes,
// passing a summary line via stdin or as an argument.
type ExecHandler struct {
	cmd  string
	args []string
}

// NewExecHandler creates an ExecHandler that invokes cmd with the
// provided args. The formatted change summary is appended as the
// final argument at call time.
func NewExecHandler(cmd string, args ...string) *ExecHandler {
	return &ExecHandler{cmd: cmd, args: args}
}

// Handle runs the configured command once per non-empty change set.
func (h *ExecHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	msg := formatExecMsg(changes)
	fullArgs := append(h.args, msg)
	c := exec.Command(h.cmd, fullArgs...) //nolint:gosec
	var stderr bytes.Buffer
	c.Stderr = &stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("exec handler: command %q failed: %w (stderr: %s)",
			h.cmd, err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// Drain is a no-op for ExecHandler.
func (h *ExecHandler) Drain() error { return nil }

func formatExecMsg(changes []monitor.Change) string {
	parts := make([]string, 0, len(changes))
	for _, c := range changes {
		action := "added"
		if c.Kind == monitor.Removed {
			action = "removed"
		}
		host := c.Binding.Hostname
		if host == "" {
			host = c.Binding.Addr
		}
		parts = append(parts, fmt.Sprintf("%s:%d/%s %s",
			host, c.Binding.Port, strings.ToLower(c.Binding.Proto), action))
	}
	return strings.Join(parts, "; ")
}
