package config

import "fmt"

// DatadogConfig holds settings for the Datadog alert handler.
type DatadogConfig struct {
	Enabled bool     `toml:"enabled" yaml:"enabled"`
	APIKey  string   `toml:"api_key" yaml:"api_key"`
	Site    string   `toml:"site"    yaml:"site"`
	Tags    []string `toml:"tags"    yaml:"tags"`
}

// DefaultDatadogConfig returns a DatadogConfig with sensible defaults.
func DefaultDatadogConfig() DatadogConfig {
	return DatadogConfig{
		Enabled: false,
		Site:    "datadoghq.com",
		Tags:    []string{},
	}
}

// Validate returns an error if the DatadogConfig is invalid.
func (c DatadogConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.APIKey == "" {
		return fmt.Errorf("datadog: api_key is required when enabled")
	}
	if c.Site == "" {
		return fmt.Errorf("datadog: site must not be empty")
	}
	return nil
}
