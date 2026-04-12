package config

import "fmt"

// DiscordConfig holds configuration for the Discord alert handler.
type DiscordConfig struct {
	// Enabled controls whether Discord alerts are sent.
	Enabled bool `toml:"enabled" json:"enabled"`

	// WebhookURL is the Discord incoming webhook URL.
	WebhookURL string `toml:"webhook_url" json:"webhook_url"`
}

// DefaultDiscordConfig returns a DiscordConfig with safe defaults.
func DefaultDiscordConfig() DiscordConfig {
	return DiscordConfig{
		Enabled:    false,
		WebhookURL: "",
	}
}

// Validate returns an error if the DiscordConfig is invalid.
func (c DiscordConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.WebhookURL == "" {
		return fmt.Errorf("discord: webhook_url must be set when enabled")
	}
	return nil
}
