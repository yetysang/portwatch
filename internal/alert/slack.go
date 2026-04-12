package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// SlackHandler sends alert notifications to a Slack incoming webhook URL.
type SlackHandler struct {
	webhookURL string
	client     *http.Client
	prefix     string
}

type slackPayload struct {
	Text string `json:"text"`
}

// NewSlackHandler creates a SlackHandler that posts messages to the given
// Slack incoming webhook URL. prefix is prepended to each message.
func NewSlackHandler(webhookURL, prefix string) *SlackHandler {
	return &SlackHandler{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
		prefix:     prefix,
	}
}

// Handle posts a Slack message for each change in the batch.
// Returns the first error encountered, if any.
func (h *SlackHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		msg := formatSlackMsg(c, h.prefix)
		if err := h.post(msg); err != nil {
			return err
		}
	}
	return nil
}

// Drain is a no-op for SlackHandler.
func (h *SlackHandler) Drain() error { return nil }

func (h *SlackHandler) post(text string) error {
	payload := slackPayload{Text: text}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}
	resp, err := h.client.Post(h.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatSlackMsg(c monitor.Change, prefix string) string {
	host := c.Binding.Hostname
	if host == "" {
		host = c.Binding.Addr
	}
	proto := c.Binding.Proto
	if proto == "" {
		proto = "tcp"
	}
	action := "added"
	if c.Kind == monitor.Removed {
		action = "removed"
	}
	process := c.Binding.Process
	if process == "" {
		process = "unknown"
	}
	msg := fmt.Sprintf("%s port %s/%d (%s) %s [pid %d, %s]",
		action, proto, c.Binding.Port, host, action, c.Binding.PID, process)
	if prefix != "" {
		return prefix + " " + msg
	}
	return msg
}
