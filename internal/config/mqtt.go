package config

import (
	"fmt"
	"time"
)

// MQTTConfig holds configuration for the MQTT alert handler.
type MQTTConfig struct {
	Enabled  bool          `toml:"enabled" yaml:"enabled"`
	Broker   string        `toml:"broker" yaml:"broker"`
	Topic    string        `toml:"topic" yaml:"topic"`
	ClientID string        `toml:"client_id" yaml:"client_id"`
	Username string        `toml:"username" yaml:"username"`
	Password string        `toml:"password" yaml:"password"`
	QoS      byte          `toml:"qos" yaml:"qos"`
	Retain   bool          `toml:"retain" yaml:"retain"`
	Timeout  time.Duration `toml:"timeout" yaml:"timeout"`
}

// DefaultMQTTConfig returns a MQTTConfig with sensible defaults.
func DefaultMQTTConfig() MQTTConfig {
	return MQTTConfig{
		Enabled:  false,
		Broker:   "tcp://localhost:1883",
		Topic:    "portwatch/alerts",
		ClientID: "portwatch",
		QoS:      0,
		Retain:   false,
		Timeout:  5 * time.Second,
	}
}

// Validate checks MQTTConfig for required fields when enabled.
func (c MQTTConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Broker == "" {
		return fmt.Errorf("mqtt: broker URL is required when enabled")
	}
	if c.Topic == "" {
		return fmt.Errorf("mqtt: topic is required when enabled")
	}
	if c.QoS > 2 {
		return fmt.Errorf("mqtt: qos must be 0, 1, or 2, got %d", c.QoS)
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("mqtt: timeout must be positive")
	}
	return nil
}
