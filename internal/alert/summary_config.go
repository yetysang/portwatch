package alert

import (
	"errors"
	"time"
)

// SummaryConfig holds configuration for the SummaryHandler.
type SummaryConfig struct {
	// Enabled controls whether the summary handler is active.
	Enabled bool `yaml:"enabled" json:"enabled"`

	// Interval is how often the accumulated summary is flushed.
	// Accepts Go duration strings, e.g. "5m", "1h".
	Interval time.Duration `yaml:"interval" json:"interval"`

	// Prefix is an optional string prepended to each summary line.
	Prefix string `yaml:"prefix" json:"prefix"`

	// OutputPath is the file to write summaries to.
	// Leave empty to write to stdout.
	OutputPath string `yaml:"output_path" json:"output_path"`
}

// DefaultSummaryConfig returns a SummaryConfig with sensible defaults.
func DefaultSummaryConfig() SummaryConfig {
	return SummaryConfig{
		Enabled:  false,
		Interval: 5 * time.Minute,
		Prefix:   "[portwatch] ",
	}
}

// Validate checks that the SummaryConfig fields are within acceptable ranges.
func (c SummaryConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Interval <= 0 {
		return errors.New("summary interval must be a positive duration")
	}
	if c.Interval < 10*time.Second {
		return errors.New("summary interval must be at least 10s")
	}
	if c.Interval > 24*time.Hour {
		return errors.New("summary interval must not exceed 24h")
	}
	return nil
}

// WithDefaults returns a new SummaryConfig with zero-value fields replaced by
// their defaults from DefaultSummaryConfig. Explicitly set values are preserved.
func (c SummaryConfig) WithDefaults() SummaryConfig {
	defaults := DefaultSummaryConfig()
	if c.Interval <= 0 {
		c.Interval = defaults.Interval
	}
	if c.Prefix == "" {
		c.Prefix = defaults.Prefix
	}
	return c
}
