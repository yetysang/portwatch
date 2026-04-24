package config

import (
	"fmt"
	"time"
)

// GraylogConfig holds configuration for the Graylog GELF HTTP handler.
type GraylogConfig struct {
	Enabled  bool          `yaml:"enabled"`
	URL      string        `yaml:"url"`
	Source   string        `yaml:"source"`
	Facility string        `yaml:"facility"`
	Timeout  time.Duration `yaml:"timeout"`
}

// DefaultGraylogConfig returns a GraylogConfig with sensible defaults.
func DefaultGraylogConfig() GraylogConfig {
	return GraylogConfig{
		Enabled:  false,
		URL:      "",
		Source:   "portwatch",
		Facility: "portwatch",
		Timeout:  5 * time.Second,
	}
}

// Validate checks the GraylogConfig for correctness.
func (c GraylogConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return fmt.Errorf("graylog: url is required when enabled")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("graylog: timeout must be positive")
	}
	if c.Timeout > 30*time.Second {
		return fmt.Errorf("graylog: timeout must not exceed 30s")
	}
	return nil
}
