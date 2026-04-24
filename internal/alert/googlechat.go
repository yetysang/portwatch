package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/wneessen/portwatch/internal/config"
	"github.com/wneessen/portwatch/internal/monitor"
)

// GoogleChatHandler sends port change alerts to a Google Chat space via an
// incoming webhook URL.
type GoogleChatHandler struct {
	cfg    config.GoogleChatConfig
	client *http.Client
}

// NewGoogleChatHandler creates a GoogleChatHandler from cfg.
func NewGoogleChatHandler(cfg config.GoogleChatConfig) *GoogleChatHandler {
	return &GoogleChatHandler{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

type googleChatPayload struct {
	Text string `json:"text"`
}

func formatGoogleChatMsg(c monitor.Change) string {
	host := c.Binding.Hostname
	if host == "" {
		host = c.Binding.Addr
	}
	action := "bound"
	if c.Kind == monitor.Removed {
		action = "unbound"
	}
	msg := fmt.Sprintf("[portwatch] Port %d/%s %s on %s",
		c.Binding.Port, c.Binding.Proto, action, host)
	if c.Binding.Process != "" {
		msg += fmt.Sprintf(" (process: %s", c.Binding.Process)
		if c.Binding.PID > 0 {
			msg += fmt.Sprintf(", pid: %d", c.Binding.PID)
		}
		msg += ")"
	}
	return msg
}

// Handle posts each change to the configured Google Chat webhook.
func (h *GoogleChatHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		payload := googleChatPayload{Text: formatGoogleChatMsg(c)}
		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("googlechat: marshal: %w", err)
		}
		resp, err := h.client.Post(h.cfg.WebhookURL, "application/json", bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("googlechat: post: %w", err)
		}
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("googlechat: unexpected status %d", resp.StatusCode)
		}
	}
	return nil
}

// Drain is a no-op for this handler.
func (h *GoogleChatHandler) Drain() error { return nil }

// ensure compile-time interface satisfaction
var _ Handler = (*GoogleChatHandler)(nil)

// keep time import used via Timeout field indirectly
var _ = time.Second
