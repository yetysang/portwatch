package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
)

// adaptiveCardHandler sends Adaptive Card payloads to an MS Teams incoming webhook.
type adaptiveCardHandler struct {
	cfg    config.AdaptiveCardConfig
	client *http.Client
}

// NewAdaptiveCardHandler creates a handler that posts Adaptive Cards to MS Teams.
func NewAdaptiveCardHandler(cfg config.AdaptiveCardConfig) Handler {
	return &adaptiveCardHandler{
		cfg: cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

func (h *adaptiveCardHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	payload := buildAdaptiveCardPayload(changes, h.cfg.ThemeColor)
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("adaptive_card: marshal payload: %w", err)
	}
	resp, err := h.client.Post(h.cfg.WebhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("adaptive_card: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("adaptive_card: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func (h *adaptiveCardHandler) Drain() error { return nil }

func buildAdaptiveCardPayload(changes []monitor.Change, themeColor string) map[string]any {
	var facts []map[string]string
	for _, c := range changes {
		action := "added"
		if c.Kind == monitor.Removed {
			action = "removed"
		}
		host := c.Binding.Hostname
		if host == "" {
			host = c.Binding.Addr
		}
		facts = append(facts, map[string]string{
			"title": fmt.Sprintf("%s:%d (%s)", host, c.Binding.Port, c.Binding.Proto),
			"value": action,
		})
	}
	return map[string]any{
		"type":       "message",
		"themeColor": themeColor,
		"summary":    fmt.Sprintf("portwatch: %d port change(s) at %s", len(changes), time.Now().UTC().Format(time.RFC3339)),
		"sections": []map[string]any{
			{
				"activityTitle":    "Port binding changes detected",
				"activitySubtitle": time.Now().UTC().Format(time.RFC3339),
				"facts":            facts,
				"markdown":         true,
			},
		},
	}
}
