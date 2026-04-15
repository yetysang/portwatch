package config

import "fmt"

// KafkaConfig holds configuration for the Kafka alert handler.
type KafkaConfig struct {
	Enabled  bool     `toml:"enabled"`
	Brokers  []string `toml:"brokers"`
	Topic    string   `toml:"topic"`
	ClientID string   `toml:"client_id"`
}

// DefaultKafkaConfig returns a KafkaConfig with sensible defaults.
func DefaultKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Enabled:  false,
		Brokers:  []string{"localhost:9092"},
		Topic:    "portwatch",
		ClientID: "portwatch",
	}
}

// Validate returns an error if the KafkaConfig is invalid.
func (c KafkaConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if len(c.Brokers) == 0 {
		return fmt.Errorf("kafka: at least one broker address is required")
	}
	for i, b := range c.Brokers {
		if b == "" {
			return fmt.Errorf("kafka: broker at index %d must not be empty", i)
		}
	}
	if c.Topic == "" {
		return fmt.Errorf("kafka: topic is required")
	}
	if c.ClientID == "" {
		return fmt.Errorf("kafka: client_id is required")
	}
	return nil
}
