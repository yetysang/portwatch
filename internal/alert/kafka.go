package alert

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/netwatch/portwatch/internal/monitor"
)

// KafkaProducer is a minimal interface for publishing messages to Kafka.
type KafkaProducer interface {
	Publish(topic string, key, value []byte) error
	Close() error
}

// KafkaHandler sends port change alerts to a Kafka topic.
type KafkaHandler struct {
	producer KafkaProducer
	topic    string
}

// NewKafkaHandler returns a KafkaHandler that publishes change events to the
// given Kafka topic using the provided producer.
func NewKafkaHandler(producer KafkaProducer, topic string) *KafkaHandler {
	return &KafkaHandler{producer: producer, topic: topic}
}

type kafkaEvent struct {
	Timestamp string `json:"timestamp"`
	Kind      string `json:"kind"`
	Proto     string `json:"proto"`
	Addr      string `json:"addr"`
	Port      uint16 `json:"port"`
	Process   string `json:"process,omitempty"`
	PID       int    `json:"pid,omitempty"`
}

// Handle publishes each change as a JSON message to the configured Kafka topic.
func (h *KafkaHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, c := range changes {
		ev := kafkaEvent{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Kind:      string(c.Kind),
			Proto:     c.Binding.Proto,
			Addr:      c.Binding.Addr,
			Port:      c.Binding.Port,
			Process:   c.Binding.Process,
			PID:       c.Binding.PID,
		}
		value, err := json.Marshal(ev)
		if err != nil {
			return fmt.Errorf("kafka: marshal event: %w", err)
		}
		key := []byte(fmt.Sprintf("%s:%d", ev.Proto, ev.Port))
		if err := h.producer.Publish(h.topic, key, value); err != nil {
			return fmt.Errorf("kafka: publish to %q: %w", h.topic, err)
		}
	}
	return nil
}

// Drain is a no-op for the Kafka handler.
func (h *KafkaHandler) Drain() error { return nil }
