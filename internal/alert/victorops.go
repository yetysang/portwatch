package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"portwatch/internal/monitor"
)

// VictorOpsHandler sends alert notifications to a VictorOps REST endpoint.
type VictorOpsHandler struct {
	webhookURL string
	client     *http.Client
}

// NewVictorOpsHandler creates a handler that posts to the given VictorOps
// REST endpoint URL (typically includes routing key).
func NewVictorOpsHandler(webhookURL string) *VictorOpsHandler {
	return &VictorOpsHandler{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Handle sends a VictorOps alert for each non-empty batch of changes.
func (h *VictorOpsHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	payload := buildVictorOpsPayload(changes[0])
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("victorops: marshal payload: %w", err)
	}
	resp, err := h.client.Post(h.webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("victorops: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("victorops: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// Drain is a no-op for VictorOpsHandler.
func (h *VictorOpsHandler) Drain() error { return nil }

func buildVictorOpsPayload(c monitor.Change) map[string]interface{} {
	msgType := "CRITICAL"
	if c.Kind == monitor.ChangeRemoved {
		msgType = "RECOVERY"
	}
	return map[string]interface{}{
		"message_type":    msgType,
		"entity_id":       fmt.Sprintf("%s/%d", c.Binding.Proto, c.Binding.Port),
		"state_message":   formatVictorOpsMsg(c),
		"monitoring_tool": "portwatch",
	}
}

func formatVictorOpsMsg(c monitor.Change) string {
	action := "bound to"
	if c.Kind == monitor.ChangeRemoved {
		action = "released from"
	}
	host := c.Binding.Hostname
	if host == "" {
		host = c.Binding.Addr
	}
	return fmt.Sprintf("port %d/%s %s %s", c.Binding.Port, c.Binding.Proto, action, host)
}
