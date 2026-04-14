package alert

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/iamcalledned/portwatch/internal/monitor"
)

// AMQPPublisher is a minimal interface for publishing to an AMQP exchange.
type AMQPPublisher interface {
	Publish(exchange, routingKey string, body []byte) error
	Close() error
}

// AMQPHandler sends port change alerts to an AMQP broker (e.g. RabbitMQ).
type AMQPHandler struct {
	publisher  AMQPPublisher
	exchange   string
	routingKey string
}

// NewAMQPHandler constructs an AMQPHandler with the given publisher and routing config.
func NewAMQPHandler(publisher AMQPPublisher, exchange, routingKey string) *AMQPHandler {
	return &AMQPHandler{
		publisher:  publisher,
		exchange:   exchange,
		routingKey: routingKey,
	}
}

type amqpMessage struct {
	Timestamp string          `json:"timestamp"`
	Action    string          `json:"action"`
	Proto     string          `json:"proto"`
	Addr      string          `json:"addr"`
	Port      uint16          `json:"port"`
	Process   string          `json:"process,omitempty"`
	PID       int             `json:"pid,omitempty"`
}

func formatAMQPMsg(c monitor.Change) amqpMessage {
	addr := c.Binding.Hostname
	if addr == "" {
		addr = c.Binding.Addr
	}
	return amqpMessage{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Action:    string(c.Kind),
		Proto:     c.Binding.Proto,
		Addr:      addr,
		Port:      c.Binding.Port,
		Process:   c.Binding.Process,
		PID:       c.Binding.PID,
	}
}

// Handle publishes each change as a JSON message to the configured AMQP exchange.
func (h *AMQPHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		msg := formatAMQPMsg(c)
		body, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("amqp: marshal: %w", err)
		}
		if err := h.publisher.Publish(h.exchange, h.routingKey, body); err != nil {
			return fmt.Errorf("amqp: publish: %w", err)
		}
	}
	return nil
}

// Drain is a no-op for AMQP.
func (h *AMQPHandler) Drain() error { return nil }
