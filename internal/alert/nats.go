package alert

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"

	"portwatch/internal/monitor"
)

// natsPublisher abstracts nats.Conn for testing.
type natsPublisher interface {
	Publish(subject string, data []byte) error
	Drain() error
}

// NatsHandler publishes port change events to a NATS subject.
type NatsHandler struct {
	conn    natsPublisher
	subject string
}

// NatsPayload is the JSON structure published to NATS.
type NatsPayload struct {
	Timestamp string `json:"timestamp"`
	Kind      string `json:"kind"`
	Proto     string `json:"proto"`
	Addr      string `json:"addr"`
	Port      int    `json:"port"`
	Process   string `json:"process,omitempty"`
	PID       int    `json:"pid,omitempty"`
}

// NewNatsHandler creates a NatsHandler using the provided nats.Conn and subject.
func NewNatsHandler(conn *nats.Conn, subject string) *NatsHandler {
	return &NatsHandler{conn: conn, subject: subject}
}

// Handle publishes each change as a JSON message to the configured NATS subject.
func (h *NatsHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		payload := NatsPayload{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Kind:      string(c.Kind),
			Proto:     c.Binding.Proto,
			Addr:      formatAddr(c.Binding),
			Port:      c.Binding.Port,
			Process:   c.Binding.Process,
			PID:       c.Binding.PID,
		}
		data, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("nats: marshal: %w", err)
		}
		if err := h.conn.Publish(h.subject, data); err != nil {
			return fmt.Errorf("nats: publish: %w", err)
		}
	}
	return nil
}

// Drain flushes pending messages and closes the connection gracefully.
func (h *NatsHandler) Drain() error {
	return h.conn.Drain()
}

func formatAddr(b interface{ GetAddr() string }) string {
	type addrGetter interface {
		GetAddr() string
	}
	if ag, ok := b.(addrGetter); ok {
		return ag.GetAddr()
	}
	return ""
}
