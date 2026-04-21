package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/wolveix/portwatch/internal/config"
	"github.com/wolveix/portwatch/internal/monitor"
)

// MattermostHandler sends alerts to a Mattermost incoming webhook.
type MattermostHandler struct {
	cfg    config.MattermostConfig
	client *http.Client
}

// NewMattermostHandler creates a MattermostHandler from cfg.
func NewMattermostHandler(cfg config.MattermostConfig) *MattermostHandler {
	return &MattermostHandler{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

type mattermostPayload struct {
	Text     string `json:"text"`
	Channel  string `json:"channel"`
	Username string `json:"username,omitempty"`
	IconURL  string `json:"icon_url,omitempty"`
}

// Handle posts each change to the configured Mattermost webhook.
func (h *MattermostHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		msg := formatMattermostMsg(c)
		payload := mattermostPayload{
			Text:     msg,
			Channel:  h.cfg.Channel,
			Username: h.cfg.Username,
			IconURL:  h.cfg.IconURL,
		}
		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("mattermost: marshal payload: %w", err)
		}
		resp, err := h.client.Post(h.cfg.URL, "application/json", bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("mattermost: post: %w", err)
		}
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("mattermost: unexpected status %d", resp.StatusCode)
		}
	}
	return nil
}

// Drain is a no-op for this handler.
func (h *MattermostHandler) Drain() error { return nil }

func formatMattermostMsg(c monitor.Change) string {
	action := "bound"
	if c.Kind == monitor.Removed {
		action = "unbound"
	}
	host := c.Binding.Hostname
	if host == "" {
		host = c.Binding.Addr
	}
	return fmt.Sprintf("**Port %s %s** — %s:%d (%s) at %s",
		action, c.Binding.Proto,
		host, c.Binding.Port,
		c.Binding.Process,
		time.Now().UTC().Format(time.RFC3339),
	)
}
