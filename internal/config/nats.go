package config

import "fmt"

// NatsConfig holds configuration for the NATS alert handler.
type NatsConfig struct {
	Enabled  bool   `toml:"enabled"`
	URL      string `toml:"url"`
	Subject  string `toml:"subject"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

// DefaultNatsConfig returns a NatsConfig with sensible defaults.
func DefaultNatsConfig() NatsConfig {
	return NatsConfig{
		Enabled:  false,
		URL:      "nats://localhost:4222",
		Subject:  "portwatch.events",
		Username: "",
		Password: "",
	}
}

// Validate checks that the NatsConfig is consistent.
func (c NatsConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return fmt.Errorf("nats: url is required when enabled")
	}
	if c.Subject == "" {
		return fmt.Errorf("nats: subject is required when enabled")
	}
	return nil
}
