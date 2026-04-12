package config

import "fmt"

// TeamsConfig holds configuration for the Microsoft Teams alert handler.
type TeamsConfig struct {
	// Enabled controls whether the Teams handler is active.
	Enabled bool `toml:"enabled" yaml:"enabled"`

	// WebhookURL is the Incoming Webhook URL for the target Teams channel.
	WebhookURL string `toml:"webhook_url" yaml:"webhook_url"`
}

// DefaultTeamsConfig returns a TeamsConfig with safe defaults (disabled).
func DefaultTeamsConfig() TeamsConfig {
	return TeamsConfig{
		Enabled:    false,
		WebhookURL: "",
	}
}

// Validate checks that the TeamsConfig is internally consistent.
// If Enabled is false, no further checks are performed.
func (c TeamsConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.WebhookURL == "" {
		return fmt.Errorf("teams: webhook_url must not be empty when enabled")
	}
	return nil
}
