package config

import (
	"fmt"
	"time"
)

// NewRelicConfig holds configuration for the New Relic alert handler.
type NewRelicConfig struct {
	Enabled   bool          `yaml:"enabled"`
	APIKey    string        `yaml:"api_key"`
	AccountID string        `yaml:"account_id"`
	Region    string        `yaml:"region"` // "US" or "EU"
	Timeout   time.Duration `yaml:"timeout"`
}

// DefaultNewRelicConfig returns a NewRelicConfig with sensible defaults.
func DefaultNewRelicConfig() NewRelicConfig {
	return NewRelicConfig{
		Enabled:   false,
		Region:    "US",
		Timeout:   5 * time.Second,
	}
}

// Validate returns an error if the configuration is invalid.
func (c NewRelicConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.APIKey == "" {
		return fmt.Errorf("newrelic: api_key is required when enabled")
	}
	if c.AccountID == "" {
		return fmt.Errorf("newrelic: account_id is required when enabled")
	}
	if c.Region != "US" && c.Region != "EU" {
		return fmt.Errorf("newrelic: region must be \"US\" or \"EU\", got %q", c.Region)
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("newrelic: timeout must be positive")
	}
	return nil
}
