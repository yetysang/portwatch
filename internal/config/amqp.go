package config

import "fmt"

// AMQPConfig holds configuration for the AMQP alert handler.
type AMQPConfig struct {
	Enabled    bool   `toml:"enabled" json:"enabled"`
	URL        string `toml:"url" json:"url"`
	Exchange   string `toml:"exchange" json:"exchange"`
	RoutingKey string `toml:"routing_key" json:"routing_key"`
	VHost      string `toml:"vhost" json:"vhost"`
}

// DefaultAMQPConfig returns an AMQPConfig with sensible defaults.
func DefaultAMQPConfig() AMQPConfig {
	return AMQPConfig{
		Enabled:    false,
		URL:        "amqp://guest:guest@localhost:5672/",
		Exchange:   "portwatch",
		RoutingKey: "port.events",
		VHost:      "/",
	}
}

// Validate returns an error if the AMQPConfig is enabled but misconfigured.
func (c AMQPConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return fmt.Errorf("amqp: url is required when enabled")
	}
	if c.Exchange == "" {
		return fmt.Errorf("amqp: exchange is required when enabled")
	}
	if c.RoutingKey == "" {
		return fmt.Errorf("amqp: routing_key is required when enabled")
	}
	return nil
}
