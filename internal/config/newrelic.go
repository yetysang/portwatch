package config

import (
	"fmt"
	"time"
)

// NewRelicConfig holds configuration for the New Relic alert handler.
type NewRelicConfig struct {
	Enabled    bool          `toml:"enabled"`
	APIKey     string        `toml:"api_key"`
	AccountID  string        `toml:"account_id"`
	EventType  string        `toml:"event_type"`
	Region     string        `toml:"region"` // "us" or "eu"
	Timeout    time.Duration `toml:"timeout"`
}

// DefaultNewRelicConfig returns a NewRelicConfig with sensible defaults.
func DefaultNewRelicConfig() NewRelicConfig {
	return NewRelicConfig{
		Enabled:   false,
		EventType: "PortWatchEvent",
		Region:    "us",
		Timeout:   10 * time.Second,
	}
}

// Validate checks that the NewRelicConfig is self-consistent.
func (c NewRelicConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.APIKey == "" {
		return ValidationError{Field: "api_key", Msg: "required when new_relic is enabled"}
	}
	if c.AccountID == "" {
		return ValidationError{Field: "account_id", Msg: "required when new_relic is enabled"}
	}
	if c.Region != "us" && c.Region != "eu" {
		return ValidationError{Field: "region", Msg: fmt.Sprintf("must be \"us\" or \"eu\", got %q", c.Region)}
	}
	if c.Timeout <= 0 {
		return ValidationError{Field: "timeout", Msg: "must be positive"}
	}
	return nil
}
