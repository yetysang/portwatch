package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/monitor"
)

// NtfyConfig holds configuration for the ntfy.sh push notification handler.
type NtfyConfig struct {
	Enabled  bool   `yaml:"enabled"`
	ServerURL string `yaml:"server_url"`
	Topic    string `yaml:"topic"`
	Token    string `yaml:"token,omitempty"`
	Priority string `yaml:"priority"`
}

type ntfyPayload struct {
	Topic    string   `json:"topic"`
	Title    string   `json:"title"`
	Message  string   `json:"message"`
	Priority string   `json:"priority"`
	Tags     []string `json:"tags"`
}

type ntfyHandler struct {
	cfg    NtfyConfig
	client *http.Client
}

// NewNtfyHandler creates an alert handler that posts to a ntfy.sh server.
func NewNtfyHandler(cfg NtfyConfig) Handler {
	return &ntfyHandler{cfg: cfg, client: &http.Client{}}
}

func (h *ntfyHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		if err := h.sendChange(c); err != nil {
			return err
		}
	}
	return nil
}

// sendChange posts a single change event to the ntfy server.
func (h *ntfyHandler) sendChange(c monitor.Change) error {
	payload := formatNtfyMsg(h.cfg.Topic, c)
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ntfy: marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, h.cfg.ServerURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("ntfy: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if h.cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+h.cfg.Token)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("ntfy: post: %w", err)
	}
	resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func (h *ntfyHandler) Drain() error { return nil }

func formatNtfyMsg(topic string, c monitor.Change) ntfyPayload {
	action := "added"
	tag := "white_check_mark"
	priority := "default"
	if c.Kind == monitor.Removed {
		action = "removed"
		tag = "warning"
		priority = "high"
	}
	host := c.Binding.Hostname
	if host == "" {
		host = c.Binding.Addr
	}
	msg := fmt.Sprintf("Port %s/%s on %s was %s", c.Binding.Port, c.Binding.Proto, host, action)
	if c.Binding.Process != "" {
		msg += fmt.Sprintf(" by %s (pid %d)", c.Binding.Process, c.Binding.PID)
	}
	return ntfyPayload{
		Topic:    topic,
		Title:    fmt.Sprintf("portwatch: port %s %s", c.Binding.Port, action),
		Message:  msg,
		Priority: priority,
		Tags:     []string{tag, "portwatch"},
	}
}
