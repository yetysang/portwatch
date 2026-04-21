package config

import (
	"fmt"
	"time"
)

// StatusPageConfig holds configuration for the Statuspage.io alert handler.
type StatusPageConfig struct {
	Enabled    bool          `yaml:"enabled"`
	APIKey     string        `yaml:"api_key"`
	PageID     string        `yaml:"page_id"`
	ComponentID string       `yaml:"component_id"`
	Timeout    time.Duration `yaml:"timeout"`
}

// DefaultStatusPageConfig returns a StatusPageConfig with sensible defaults.
func DefaultStatusPageConfig() StatusPageConfig {
	return StatusPageConfig{
		Enabled: false,
		Timeout: 10 * time.Second,
	}
}

// Validate returns an error if the configuration is invalid.
func (c StatusPageConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.APIKey == "" {
		return fmt.Errorf("statuspage: api_key is required when enabled")
	}
	if c.PageID == "" {
		return fmt.Errorf("statuspage: page_id is required when enabled")
	}
	if c.ComponentID == "" {
		return fmt.Errorf("statuspage: component_id is required when enabled")
	}
	if c.Timeout < time.Second {
		return fmt.Errorf("statuspage: timeout must be at least 1s")
	}
	return nil
}
