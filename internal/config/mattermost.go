package config

import (
	"fmt"
	"time"
)

// MattermostConfig holds settings for the Mattermost alert handler.
type MattermostConfig struct {
	Enabled  bool          `toml:"enabled"`
	URL      string        `toml:"url"`
	Channel  string        `toml:"channel"`
	Username string        `toml:"username"`
	IconURL  string        `toml:"icon_url"`
	Timeout  time.Duration `toml:"timeout"`
}

// DefaultMattermostConfig returns a MattermostConfig with sensible defaults.
func DefaultMattermostConfig() MattermostConfig {
	return MattermostConfig{
		Enabled:  false,
		URL:      "",
		Channel:  "",
		Username: "portwatch",
		IconURL:  "",
		Timeout:  10 * time.Second,
	}
}

// Validate returns an error if the configuration is invalid.
func (c MattermostConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return fmt.Errorf("mattermost: url is required when enabled")
	}
	if c.Channel == "" {
		return fmt.Errorf("mattermost: channel is required when enabled")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("mattermost: timeout must be positive")
	}
	return nil
}
