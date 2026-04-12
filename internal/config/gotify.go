package config

import "fmt"

// GotifyConfig holds Gotify-specific configuration values.
type GotifyConfig struct {
	Enabled  bool   `toml:"enabled"`
	URL      string `toml:"url"`
	Token    string `toml:"token"`
	Priority int    `toml:"priority"`
}

// DefaultGotifyConfig returns a GotifyConfig with sensible defaults.
func DefaultGotifyConfig() GotifyConfig {
	return GotifyConfig{
		Enabled:  false,
		URL:      "",
		Token:    "",
		Priority: 5,
	}
}

// Validate checks that the GotifyConfig is consistent.
func (c GotifyConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return fmt.Errorf("gotify: url is required when enabled")
	}
	if c.Token == "" {
		return fmt.Errorf("gotify: token is required when enabled")
	}
	if c.Priority < 0 || c.Priority > 10 {
		return fmt.Errorf("gotify: priority must be between 0 and 10, got %d", c.Priority)
	}
	return nil
}
