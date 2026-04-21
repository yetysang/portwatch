package config

import (
	"fmt"
	"time"
)

// DeadManConfig holds settings for the dead-man's switch alert handler.
// When enabled, portwatch expects to emit at least one heartbeat within
// every Interval. If no scan completes in time, the switch fires.
type DeadManConfig struct {
	Enabled  bool          `toml:"enabled"`
	Interval time.Duration `toml:"interval"`
	URL      string        `toml:"url"`   // endpoint to ping on each heartbeat
	Timeout  time.Duration `toml:"timeout"` // HTTP request timeout
}

// DefaultDeadManConfig returns a DeadManConfig with safe defaults.
func DefaultDeadManConfig() DeadManConfig {
	return DeadManConfig{
		Enabled:  false,
		Interval: 60 * time.Second,
		Timeout:  5 * time.Second,
	}
}

// Validate returns an error if the configuration is invalid.
func (c DeadManConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return fmt.Errorf("deadman: url must not be empty when enabled")
	}
	if c.Interval < 10*time.Second {
		return fmt.Errorf("deadman: interval must be at least 10s, got %s", c.Interval)
	}
	if c.Interval > 24*time.Hour {
		return fmt.Errorf("deadman: interval must not exceed 24h, got %s", c.Interval)
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("deadman: timeout must be positive, got %s", c.Timeout)
	}
	return nil
}
