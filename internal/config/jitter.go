package config

import (
	"fmt"
	"time"
)

// JitterConfig controls the random jitter applied to the polling interval
// to avoid thundering-herd effects when multiple instances run simultaneously.
type JitterConfig struct {
	// Enabled controls whether jitter is applied.
	Enabled bool `toml:"enabled" yaml:"enabled"`

	// MaxJitter is the maximum duration of random jitter added to each tick.
	// Must be less than the main poll interval.
	MaxJitter time.Duration `toml:"max_jitter" yaml:"max_jitter"`
}

// DefaultJitterConfig returns a JitterConfig with sensible defaults.
func DefaultJitterConfig() JitterConfig {
	return JitterConfig{
		Enabled:   false,
		MaxJitter: 500 * time.Millisecond,
	}
}

// Validate checks that the JitterConfig is self-consistent.
func (c JitterConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.MaxJitter <= 0 {
		return fmt.Errorf("jitter: max_jitter must be positive, got %s", c.MaxJitter)
	}
	if c.MaxJitter > 30*time.Second {
		return fmt.Errorf("jitter: max_jitter %s exceeds maximum of 30s", c.MaxJitter)
	}
	return nil
}
