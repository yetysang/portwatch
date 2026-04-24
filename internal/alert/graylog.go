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

// graylogPayload represents a GELF message sent to Graylog over HTTP.
type graylogPayload struct {
	Version      string  `json:"version"`
	Host         string  `json:"host"`
	ShortMessage string  `json:"short_message"`
	Timestamp    float64 `json:"timestamp"`
	Level        int     `json:"level"`
	Facility     string  `json:"facility"`
	Action       string  `json:"_action"`
	Proto        string  `json:"_proto"`
	Port         int     `json:"_port"`
	Process      string  `json:"_process"`
}

// GraylogHandler sends GELF messages to a Graylog HTTP input.
type GraylogHandler struct {
	cfg    config.GraylogConfig
	client *http.Client
}

// NewGraylogHandler creates a new GraylogHandler from the provided config.
func NewGraylogHandler(cfg config.GraylogConfig) *GraylogHandler {
	return &GraylogHandler{
		cfg: cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

// Handle sends each change as a GELF message to the configured Graylog endpoint.
func (h *GraylogHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		if err := h.send(c); err != nil {
			return err
		}
	}
	return nil
}

// Drain is a no-op for GraylogHandler.
func (h *GraylogHandler) Drain() error { return nil }

func (h *GraylogHandler) send(c monitor.Change) error {
	payload := formatGraylogMsg(c, h.cfg)
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("graylog: marshal: %w", err)
	}
	resp, err := h.client.Post(h.cfg.URL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("graylog: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("graylog: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatGraylogMsg(c monitor.Change, cfg config.GraylogConfig) graylogPayload {
	host := c.Binding.Hostname
	if host == "" {
		host = c.Binding.Addr
	}
	level := 6 // informational
	if c.Kind == monitor.ChangeAdded {
		level = 5 // notice
	}
	return graylogPayload{
		Version:      "1.1",
		Host:         cfg.Source,
		ShortMessage: fmt.Sprintf("%s %s:%d/%s", c.Kind, host, c.Binding.Port, c.Binding.Proto),
		Timestamp:    float64(time.Now().UnixNano()) / 1e9,
		Level:        level,
		Facility:     cfg.Facility,
		Action:       string(c.Kind),
		Proto:        c.Binding.Proto,
		Port:         c.Binding.Port,
		Process:      c.Binding.Process,
	}
}
