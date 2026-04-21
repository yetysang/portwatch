package config

import (
	"fmt"
	"time"
)

// RateLimitConfig controls per-port alert rate limiting.
type RateLimitConfig struct {
	// Enabled toggles rate limiting globally.
	Enabled bool `toml:"enabled" yaml:"enabled"`

	// Cooldown is the minimum duration between repeated alerts for the same
	// port+protocol+kind triple.
	Cooldown time.Duration `toml:"cooldown" yaml:"cooldown"`

	// MaxBurst is the number of alerts allowed before the cooldown kicks in.
	MaxBurst int `toml:"max_burst" yaml:"max_burst"`
}

// DefaultRateLimitConfig returns a RateLimitConfig with sensible defaults.
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Enabled:  true,
		Cooldown: 30 * time.Second,
		MaxBurst: 1,
	}
}

// Validate returns an error if the configuration is invalid.
func (c RateLimitConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Cooldown < time.Second {
		return fmt.Errorf("ratelimit: cooldown must be at least 1s, got %s", c.Cooldown)
	}
	if c.Cooldown > 24*time.Hour {
		return fmt.Errorf("ratelimit: cooldown must not exceed 24h, got %s", c.Cooldown)
	}
	if c.MaxBurst < 1 {
		return fmt.Errorf("ratelimit: max_burst must be at least 1, got %d", c.MaxBurst)
	}
	if c.MaxBurst > 100 {
		return fmt.Errorf("ratelimit: max_burst must not exceed 100, got %d", c.MaxBurst)
	}
	return nil
}
