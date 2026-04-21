package config

import (
	"fmt"
	"time"
)

// ZScoreConfig holds configuration for the anomaly detection handler
// that uses a rolling Z-score to flag statistically unusual port churn.
type ZScoreConfig struct {
	Enabled     bool          `toml:"enabled" json:"enabled"`
	WindowSize  int           `toml:"window_size" json:"window_size"`
	Threshold   float64       `toml:"threshold" json:"threshold"`
	MinSamples  int           `toml:"min_samples" json:"min_samples"`
	Cooldown    time.Duration `toml:"cooldown" json:"cooldown"`
}

// DefaultZScoreConfig returns a ZScoreConfig with sensible defaults.
func DefaultZScoreConfig() ZScoreConfig {
	return ZScoreConfig{
		Enabled:    false,
		WindowSize: 60,
		Threshold:  3.0,
		MinSamples: 10,
		Cooldown:   5 * time.Minute,
	}
}

// Validate checks that the ZScoreConfig fields are within acceptable ranges.
func (z ZScoreConfig) Validate() error {
	if !z.Enabled {
		return nil
	}
	if z.WindowSize < 5 || z.WindowSize > 1000 {
		return validationErrorf("zscore.window_size must be between 5 and 1000, got %d", z.WindowSize)
	}
	if z.Threshold < 1.0 || z.Threshold > 10.0 {
		return validationErrorf("zscore.threshold must be between 1.0 and 10.0, got %g", z.Threshold)
	}
	if z.MinSamples < 2 || z.MinSamples > z.WindowSize {
		return fmt.Errorf("zscore.min_samples must be between 2 and window_size (%d), got %d", z.WindowSize, z.MinSamples)
	}
	if z.Cooldown < 10*time.Second || z.Cooldown > 24*time.Hour {
		return validationErrorf("zscore.cooldown must be between 10s and 24h, got %s", z.Cooldown)
	}
	return nil
}
