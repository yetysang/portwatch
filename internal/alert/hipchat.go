package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/warden-protocol/portwatch/internal/config"
	"github.com/warden-protocol/portwatch/internal/monitor"
)

// HipChatHandler sends port change alerts to a HipChat room.
type HipChatHandler struct {
	cfg    config.HipChatConfig
	client *http.Client
}

// NewHipChatHandler creates a HipChatHandler from the given config.
func NewHipChatHandler(cfg config.HipChatConfig) *HipChatHandler {
	return &HipChatHandler{
		cfg: cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

type hipchatPayload struct {
	Message       string `json:"message"`
	MessageFormat string `json:"message_format"`
	Color         string `json:"color"`
	Notify        bool   `json:"notify"`
}

// Handle sends a HipChat notification for each change.
func (h *HipChatHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		if err := h.post(formatHipChatMsg(c)); err != nil {
			return err
		}
	}
	return nil
}

// Drain is a no-op for HipChatHandler.
func (h *HipChatHandler) Drain() error { return nil }

func (h *HipChatHandler) post(msg string) error {
	payload := hipchatPayload{
		Message:       msg,
		MessageFormat: "text",
		Color:         h.cfg.Color,
		Notify:        h.cfg.Notify,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("hipchat: marshal payload: %w", err)
	}
	url := fmt.Sprintf("%s/v2/room/%s/notification", h.cfg.BaseURL, h.cfg.RoomID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("hipchat: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.cfg.AuthToken)
	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("hipchat: send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("hipchat: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatHipChatMsg(c monitor.Change) string {
	host := c.Binding.Hostname
	if host == "" {
		host = c.Binding.Addr
	}
	action := "bound"
	if c.Kind == monitor.Removed {
		action = "unbound"
	}
	proc := c.Binding.Process
	if proc == "" {
		proc = "unknown"
	}
	return fmt.Sprintf("[portwatch] port %s/%d %s on %s (pid %d, %s) at %s",
		c.Binding.Proto, c.Binding.Port, action, host,
		c.Binding.PID, proc, time.Now().UTC().Format(time.RFC3339))
}
