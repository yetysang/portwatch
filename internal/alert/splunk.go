package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// SplunkHandler sends port change events to a Splunk HTTP Event Collector (HEC) endpoint.
type SplunkHandler struct {
	url    string
	token  string
	client *http.Client
}

type splunkEvent struct {
	Time   float64        `json:"time"`
	Source string         `json:"source"`
	Event  map[string]any `json:"event"`
}

// NewSplunkHandler returns a Handler that forwards changes to Splunk HEC.
func NewSplunkHandler(url, token string) *SplunkHandler {
	return &SplunkHandler{
		url:    url,
		token:  token,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (h *SplunkHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	var buf bytes.Buffer
	for _, c := range changes {
		ev := splunkEvent{
			Time:   float64(time.Now().UnixNano()) / 1e9,
			Source: "portwatch",
			Event:  formatSplunkEvent(c),
		}
		data, err := json.Marshal(ev)
		if err != nil {
			return fmt.Errorf("splunk: marshal event: %w", err)
		}
		buf.Write(data)
		buf.WriteByte('\n')
	}
	req, err := http.NewRequest(http.MethodPost, h.url, &buf)
	if err != nil {
		return fmt.Errorf("splunk: build request: %w", err)
	}
	req.Header.Set("Authorization", "Splunk "+h.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("splunk: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("splunk: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func (h *SplunkHandler) Drain() error { return nil }

func formatSplunkEvent(c monitor.Change) map[string]any {
	b := c.Binding
	host := b.Hostname
	if host == "" {
		host = b.Addr
	}
	return map[string]any{
		"action":   string(c.Kind),
		"host":     host,
		"port":     b.Port,
		"proto":    b.Proto,
		"pid":      b.PID,
		"process":  b.Process,
		"service":  b.Service,
	}
}
