package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/wrossmorrow/portwatch/internal/config"
	"github.com/wrossmorrow/portwatch/internal/monitor"
)

const zendutyAPIBase = "https://www.zenduty.com/api/events"

// ZendutyHandler sends port change alerts to Zenduty.
type ZendutyHandler struct {
	cfg    config.ZendutyConfig
	client *http.Client
}

// NewZendutyHandler creates a ZendutyHandler from the given config.
func NewZendutyHandler(cfg config.ZendutyConfig) *ZendutyHandler {
	return &ZendutyHandler{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type zendutyPayload struct {
	AlertType string `json:"alert_type"`
	Message   string `json:"message"`
	Summary   string `json:"summary"`
	EntityID  string `json:"entity_id"`
}

func formatZendutyMsg(c monitor.Change) string {
	host := c.Binding.Hostname
	if host == "" {
		host = c.Binding.Addr
	}
	action := "bound"
	if c.Kind == monitor.Removed {
		action = "unbound"
	}
	return fmt.Sprintf("port %d/%s %s on %s", c.Binding.Port, c.Binding.Proto, action, host)
}

// Handle sends each change as a Zenduty alert event.
func (h *ZendutyHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		msg := formatZendutyMsg(c)
		payload := zendutyPayload{
			AlertType: h.cfg.AlertType,
			Message:   msg,
			Summary:   msg,
			EntityID:  fmt.Sprintf("%s-%d-%s", c.Binding.Addr, c.Binding.Port, c.Binding.Proto),
		}
		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("zenduty: marshal: %w", err)
		}
		url := fmt.Sprintf("%s/%s/%s/", zendutyAPIBase, h.cfg.ServiceID, h.cfg.IntegrationID)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("zenduty: request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Token "+h.cfg.APIKey)
		resp, err := h.client.Do(req)
		if err != nil {
			return fmt.Errorf("zenduty: post: %w", err)
		}
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("zenduty: unexpected status %d", resp.StatusCode)
		}
	}
	return nil
}

// Drain is a no-op for ZendutyHandler.
func (h *ZendutyHandler) Drain() error { return nil }
