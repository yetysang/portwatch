package config

import (
	"fmt"
	"time"
)

// AdaptiveCardConfig holds configuration for the MS Teams Adaptive Card alert handler.
type AdaptiveCardConfig struct {
	Enabled    bool          `yaml:"enabled"`
	WebhookURL string        `yaml:"webhook_url"`
	Timeout    time.Duration `yaml:"timeout"`
	ThemeColor string        `yaml:"theme_color"`
}

// DefaultAdaptiveCardConfig returns a config with sensible defaults.
func DefaultAdaptiveCardConfig() AdaptiveCardConfig {
	return AdaptiveCardConfig{
		Enabled:    false,
		WebhookURL: "",
		Timeout:    10 * time.Second,
		ThemeColor: "0078D4",
	}
}

// Validate checks the AdaptiveCardConfig for correctness.
func (c AdaptiveCardConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.WebhookURL == "" {
		return fmt.Errorf("adaptive_card: webhook_url is required when enabled")
	}
	if c.Timeout < time.Second {
		return fmt.Errorf("adaptive_card: timeout must be at least 1s")
	}
	if c.Timeout > 60*time.Second {
		return fmt.Errorf("adaptive_card: timeout must not exceed 60s")
	}
	return nil
}
