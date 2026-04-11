package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// WebhookPayload is the JSON body sent to a webhook endpoint.
type WebhookPayload struct {
	Timestamp string `json:"timestamp"`
	Event     string `json:"event"`
	Proto     string `json:"proto"`
	Addr      string `json:"addr"`
	Port      int    `json:"port"`
	Process   string `json:"process,omitempty"`
	PID       int    `json:"pid,omitempty"`
}

// WebhookHandler sends change alerts to an HTTP endpoint as JSON.
type WebhookHandler struct {
	URL    string
	client *http.Client
}

// NewWebhookHandler creates a WebhookHandler that posts to the given URL.
func NewWebhookHandler(url string, timeout time.Duration) *WebhookHandler {
	return &WebhookHandler{
		URL: url,
		client: &http.Client{Timeout: timeout},
	}
}

// Handle sends each change in the batch to the configured webhook URL.
func (w *WebhookHandler) Handle(changes []monitor.Change) error {
	for _, c := range changes {
		payload := WebhookPayload{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Event:     string(c.Type),
			Proto:     c.Binding.Proto,
			Addr:      c.Binding.Addr,
			Port:      c.Binding.Port,
			Process:   c.Binding.Process,
			PID:       c.Binding.PID,
		}
		if err := w.post(payload); err != nil {
			return fmt.Errorf("webhook post failed: %w", err)
		}
	}
	return nil
}

func (w *WebhookHandler) post(payload WebhookPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := w.client.Post(w.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}
