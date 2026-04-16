package config

import "fmt"

// WebhookConfig holds configuration for the outbound webhook alert handler.
type WebhookConfig struct {
	Enabled bool   `toml:"enabled" json:"enabled"`
	URL     string `toml:"url"     json:"url"`
	Secret  string `toml:"secret"  json:"secret"`
	Timeout int    `toml:"timeout" json:"timeout"` // seconds
}

// DefaultWebhookConfig returns a WebhookConfig with sensible defaults.
func DefaultWebhookConfig() WebhookConfig {
	return WebhookConfig{
		Enabled: false,
		URL:     "",
		Secret:  "",
		Timeout: 5,
	}
}

// Validate returns an error if the WebhookConfig is invalid.
func (c WebhookConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return fmt.Errorf("webhook: url is required when enabled")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("webhook: timeout must be positive, got %d", c.Timeout)
	}
	if c.Timeout > 60 {
		return fmt.Errorf("webhook: timeout must be <= 60 seconds, got %d", c.Timeout)
	}
	return nil
}
