package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// LokiConfig holds configuration for the Grafana Loki push handler.
type LokiConfig struct {
	Enabled  bool
	URL      string // e.g. http://localhost:3100
	JobLabel string
}

type lokiHandler struct {
	cfg    LokiConfig
	client *http.Client
}

// NewLokiHandler returns a Handler that pushes log lines to a Loki instance.
func NewLokiHandler(cfg LokiConfig) Handler {
	return &lokiHandler{
		cfg:    cfg,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (h *lokiHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	streams := make([]lokiStream, 0, len(changes))
	for _, c := range changes {
		streams = append(streams, lokiStream{
			Stream: map[string]string{
				"job":   h.cfg.JobLabel,
				"proto": c.Binding.Proto,
				"kind":  string(c.Kind),
			},
			Values: [][]string{
				{strconv.FormatInt(time.Now().UnixNano(), 10), formatLokiMsg(c)},
			},
		})
	}

	body, err := json.Marshal(map[string]any{"streams": streams})
	if err != nil {
		return fmt.Errorf("loki: marshal: %w", err)
	}

	url := h.cfg.URL + "/loki/api/v1/push"
	resp, err := h.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("loki: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("loki: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func (h *lokiHandler) Drain() error { return nil }

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

func formatLokiMsg(c monitor.Change) string {
	b := &c.Binding
	host := b.Hostname
	if host == "" {
		host = b.IP
	}
	proc := b.Process
	if proc == "" {
		proc = "unknown"
	}
	return fmt.Sprintf("[portwatch] %s %s/%d addr=%s proc=%s pid=%d",
		c.Kind, b.Proto, b.Port, host, proc, b.PID)
}
