package config

import (
	"fmt"
	"time"
)

// OpsGenieConfig holds configuration for the OpsGenie alert handler.
type OpsGenieConfig struct {
	Enabled  bool          `toml:"enabled"`
	APIKey   string        `toml:"api_key"`
	Team     string        `toml:"team"`
	Priority string        `toml:"priority"`
	Timeout  time.Duration `toml:"timeout"`
}

// DefaultOpsGenieConfig returns an OpsGenieConfig with sensible defaults.
func DefaultOpsGenieConfig() OpsGenieConfig {
	return OpsGenieConfig{
		Enabled:  false,
		APIKey:   "",
		Team:     "",
		Priority: "P3",
		Timeout:  5 * time.Second,
	}
}

// Validate checks that the OpsGenieConfig is consistent.
func (c OpsGenieConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.APIKey == "" {
		return ValidationError{Field: "opsgenie.api_key", Msg: "api_key is required when opsgenie is enabled"}
	}
	validPriorities := map[string]bool{
		"P1": true, "P2": true, "P3": true, "P4": true, "P5": true,
	}
	if !validPriorities[c.Priority] {
		return ValidationError{
			Field: "opsgenie.priority",
			Msg:   fmt.Sprintf("invalid priority %q: must be one of P1–P5", c.Priority),
		}
	}
	if c.Timeout <= 0 {
		return ValidationError{Field: "opsgenie.timeout", Msg: "timeout must be positive"}
	}
	return nil
}
