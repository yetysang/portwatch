package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/wvictim/portwatch/internal/config"
	"github.com/wvictim/portwatch/internal/monitor"
)

const (
	newRelicUSEndpoint = "https://insights-collector.newrelic.com/v1/accounts/%s/events"
	newRelicEUEndpoint = "https://insights-collector.eu01.nr-data.net/v1/accounts/%s/events"
)

// newRelicEvent is the payload sent to New Relic Insights.
type newRelicEvent struct {
	EventType  string `json:"eventType"`
	Action     string `json:"action"`
	Port       int    `json:"port"`
	Proto      string `json:"proto"`
	Addr       string `json:"addr"`
	Process    string `json:"process,omitempty"`
	PID        int    `json:"pid,omitempty"`
	Timestamp  int64  `json:"timestamp"`
}

// newRelicHandler sends port change events to New Relic.
type newRelicHandler struct {
	cfg    config.NewRelicConfig
	client *http.Client
	endpoint string
}

// NewNewRelicHandler returns a Handler that posts events to New Relic Insights.
func NewNewRelicHandler(cfg config.NewRelicConfig) Handler {
	endpoint := fmt.Sprintf(newRelicUSEndpoint, cfg.AccountID)
	if cfg.Region == "EU" {
		endpoint = fmt.Sprintf(newRelicEUEndpoint, cfg.AccountID)
	}
	return &newRelicHandler{
		cfg:      cfg,
		client:   &http.Client{Timeout: cfg.Timeout},
		endpoint: endpoint,
	}
}

func (h *newRelicHandler) Handle(changes []monitor.Change) error {
	if !h.cfg.Enabled || len(changes) == 0 {
		return nil
	}
	events := make([]newRelicEvent, 0, len(changes))
	for _, c := range changes {
		events = append(events, formatNewRelicEvent(c))
	}
	body, err := json.Marshal(events)
	if err != nil {
		return fmt.Errorf("newrelic: marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, h.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("newrelic: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Insert-Key", h.cfg.APIKey)
	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("newrelic: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("newrelic: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func (h *newRelicHandler) Drain() error { return nil }

func formatNewRelicEvent(c monitor.Change) newRelicEvent {
	action := "added"
	if c.Kind == monitor.Removed {
		action = "removed"
	}
	addr := c.Binding.Addr
	if c.Binding.Hostname != "" {
		addr = c.Binding.Hostname
	}
	return newRelicEvent{
		EventType: "PortWatchEvent",
		Action:    action,
		Port:      c.Binding.Port,
		Proto:     c.Binding.Proto,
		Addr:      addr,
		Process:   c.Binding.Process,
		PID:       c.Binding.PID,
		Timestamp: time.Now().UnixMilli(),
	}
}
