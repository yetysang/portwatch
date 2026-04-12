package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

const pagerDutyEventURL = "https://events.pagerduty.com/v2/enqueue"

// PagerDutyHandler sends critical port-change alerts to PagerDuty Events API v2.
type PagerDutyHandler struct {
	routingKey string
	severity   string
	client     *http.Client
}

type pdPayload struct {
	RoutingKey  string    `json:"routing_key"`
	EventAction string    `json:"event_action"`
	Payload     pdDetails `json:"payload"`
}

type pdDetails struct {
	Summary   string `json:"summary"`
	Source    string `json:"source"`
	Severity  string `json:"severity"`
	Timestamp string `json:"timestamp"`
}

// NewPagerDutyHandler creates a handler that forwards changes to PagerDuty.
// severity should be one of: critical, error, warning, info.
func NewPagerDutyHandler(routingKey, severity string) *PagerDutyHandler {
	if severity == "" {
		severity = "error"
	}
	return &PagerDutyHandler{
		routingKey: routingKey,
		severity:   severity,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Handle sends one PagerDuty event per change.
func (h *PagerDutyHandler) Handle(changes []monitor.Change) error {
	for _, c := range changes {
		if err := h.send(c); err != nil {
			return err
		}
	}
	return nil
}

// Drain is a no-op for PagerDuty.
func (h *PagerDutyHandler) Drain() error { return nil }

func (h *PagerDutyHandler) send(c monitor.Change) error {
	summary := formatPagerDutyMsg(c)
	body := pdPayload{
		RoutingKey:  h.routingKey,
		EventAction: "trigger",
		Payload: pdDetails{
			Summary:   summary,
			Source:    "portwatch",
			Severity:  h.severity,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("pagerduty: marshal: %w", err)
	}
	resp, err := h.client.Post(pagerDutyEventURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("pagerduty: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatPagerDutyMsg(c monitor.Change) string {
	action := "added"
	if c.Kind == monitor.Removed {
		action = "removed"
	}
	host := c.Binding.Hostname
	if host == "" {
		host = c.Binding.Addr
	}
	proc := c.Binding.Process
	if proc == "" {
		proc = "unknown"
	}
	return fmt.Sprintf("port %s/%s %s on %s (pid %d, %s)",
		c.Binding.Port, c.Binding.Proto, action, host, c.Binding.PID, proc)
}
