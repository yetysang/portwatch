package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// DatadogConfig holds configuration for the Datadog Events handler.
type DatadogConfig struct {
	Enabled bool
	APIKey  string
	Site    string // e.g. "datadoghq.com" or "datadoghq.eu"
}

type datadogHandler struct {
	cfg    DatadogConfig
	client *http.Client
}

type datadogEvent struct {
	Title     string   `json:"title"`
	Text      string   `json:"text"`
	AlertType string   `json:"alert_type"` // "info", "warning", "error", "success"
	Tags      []string `json:"tags"`
	DateHappened int64 `json:"date_happened"`
}

// NewDatadogHandler returns a Handler that ships port-change events to Datadog.
func NewDatadogHandler(cfg DatadogConfig) Handler {
	if cfg.Site == "" {
		cfg.Site = "datadoghq.com"
	}
	return &datadogHandler{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (h *datadogHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		evt := formatDatadogEvent(c)
		if err := h.post(evt); err != nil {
			return err
		}
	}
	return nil
}

func (h *datadogHandler) Drain() error { return nil }

func (h *datadogHandler) post(evt datadogEvent) error {
	body, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("datadog: marshal: %w", err)
	}
	url := fmt.Sprintf("https://api.%s/api/v1/events", h.cfg.Site)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("datadog: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DD-API-KEY", h.cfg.APIKey)
	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("datadog: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("datadog: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatDatadogEvent(c monitor.Change) datadogEvent {
	action := "added"
	alertType := "warning"
	if c.Kind == monitor.Removed {
		action = "removed"
		alertType = "info"
	}
	host := c.Binding.Hostname
	if host == "" {
		host = c.Binding.Addr
	}
	title := fmt.Sprintf("portwatch: port %s/%s %s on %s",
		c.Binding.Port, c.Binding.Proto, action, host)
	text := title
	if c.Binding.Process != "" {
		text = fmt.Sprintf("%s (process: %s pid:%d)", title, c.Binding.Process, c.Binding.PID)
	}
	tags := []string{
		"source:portwatch",
		fmt.Sprintf("port:%s", c.Binding.Port),
		fmt.Sprintf("proto:%s", c.Binding.Proto),
		fmt.Sprintf("action:%s", action),
	}
	return datadogEvent{
		Title:        title,
		Text:         text,
		AlertType:    alertType,
		Tags:         tags,
		DateHappened: time.Now().Unix(),
	}
}
