package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/portwatch/internal/monitor"
)

// teamsPayload is the Adaptive Card message format for Microsoft Teams.
type teamsPayload struct {
	Type       string        `json:"type"`
	Attachments []teamsCard  `json:"attachments"`
}

type teamsCard struct {
	ContentType string      `json:"contentType"`
	Content     teamsBody   `json:"content"`
}

type teamsBody struct {
	Schema  string        `json:"$schema"`
	Type    string        `json:"type"`
	Version string        `json:"version"`
	Body    []teamsBlock  `json:"body"`
}

type teamsBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
	Wrap bool   `json:"wrap,omitempty"`
}

// TeamsHandler sends change alerts to a Microsoft Teams channel via
// an Incoming Webhook URL.
type TeamsHandler struct {
	webhookURL string
	client     *http.Client
}

// NewTeamsHandler creates a TeamsHandler that posts to the given webhook URL.
func NewTeamsHandler(webhookURL string) *TeamsHandler {
	return &TeamsHandler{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

// Handle sends a Teams message for each change in the slice.
func (h *TeamsHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		if err := h.post(formatTeamsMsg(c)); err != nil {
			return err
		}
	}
	return nil
}

// Drain is a no-op for TeamsHandler.
func (h *TeamsHandler) Drain() error { return nil }

func (h *TeamsHandler) post(msg string) error {
	payload := teamsPayload{
		Type: "message",
		Attachments: []teamsCard{
			{
				ContentType: "application/vnd.microsoft.card.adaptive",
				Content: teamsBody{
					Schema:  "http://adaptivecards.io/schemas/adaptive-card.json",
					Type:    "AdaptiveCard",
					Version: "1.4",
					Body: []teamsBlock{
						{Type: "TextBlock", Text: msg, Wrap: true},
					},
				},
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("teams: marshal payload: %w", err)
	}
	resp, err := h.client.Post(h.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("teams: http post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("teams: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatTeamsMsg(c monitor.Change) string {
	b := c.Binding
	host := b.Hostname
	if host == "" {
		host = b.IP
	}
	action := "added"
	if c.Kind == monitor.Removed {
		action = "removed"
	}
	proc := b.Process
	if proc == "" {
		proc = "unknown"
	}
	return fmt.Sprintf("[portwatch] Port %s/%s %s on %s (pid %d, process: %s)",
		b.Port, b.Proto, action, host, b.PID, proc)
}
