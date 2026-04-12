package config

import "fmt"

// SplunkConfig holds configuration for the Splunk HEC alert handler.
type SplunkConfig struct {
	Enabled  bool   `toml:"enabled"`
	URL      string `toml:"url"`
	Token    string `toml:"token"`
}

// DefaultSplunkConfig returns a SplunkConfig with safe defaults.
func DefaultSplunkConfig() SplunkConfig {
	return SplunkConfig{
		Enabled: false,
		URL:     "",
		Token:   "",
	}
}

// Validate returns an error if the SplunkConfig is enabled but incomplete.
func (c SplunkConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return fmt.Errorf("splunk: url is required when enabled")
	}
	if c.Token == "" {
		return fmt.Errorf("splunk: token is required when enabled")
	}
	return nil
}
