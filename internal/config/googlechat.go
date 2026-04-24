package config

import (
	"fmt"
	"time"
)

// GoogleChatConfig holds configuration for the Google Chat webhook alert handler.
type GoogleChatConfig struct {
	Enabled    bool          `yaml:"enabled"`
	WebhookURL string        `yaml:"webhook_url"`
	Timeout    time.Duration `yaml:"timeout"`
}

// DefaultGoogleChatConfig returns a GoogleChatConfig with safe defaults.
func DefaultGoogleChatConfig() GoogleChatConfig {
	return GoogleChatConfig{
		Enabled:    false,
		WebhookURL: "",
		Timeout:    5 * time.Second,
	}
}

// Validate checks that the GoogleChatConfig is well-formed.
func (c GoogleChatConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.WebhookURL == "" {
		return fmt.Errorf("googlechat: webhook_url is required when enabled")
	}
	if c.Timeout < time.Second {
		return fmt.Errorf("googlechat: timeout must be at least 1s")
	}
	if c.Timeout > 30*time.Second {
		return fmt.Errorf("googlechat: timeout must not exceed 30s")
	}
	return nil
}
