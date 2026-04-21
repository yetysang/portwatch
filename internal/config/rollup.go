package config

import (
	"fmt"
	"time"
)

// RollupConfig controls change-rollup (batching) behaviour before alerts are dispatched.
type RollupConfig struct {
	Enabled  bool          `toml:"enabled" yaml:"enabled"`
	Window   time.Duration `toml:"window" yaml:"window"`
	MaxBatch int           `toml:"max_batch" yaml:"max_batch"`
}

// DefaultRollupConfig returns a safe, conservative default.
func DefaultRollupConfig() RollupConfig {
	return RollupConfig{
		Enabled:  false,
		Window:   5 * time.Second,
		MaxBatch: 50,
	}
}

// Validate returns a ValidationError when the configuration is invalid.
func (c RollupConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	const (
		minWindow = 500 * time.Millisecond
		maxWindow = 5 * time.Minute
	)
	if c.Window < minWindow {
		return ValidationError{Field: "rollup.window", Msg: fmt.Sprintf("must be at least %s", minWindow)}
	}
	if c.Window > maxWindow {
		return ValidationError{Field: "rollup.window", Msg: fmt.Sprintf("must not exceed %s", maxWindow)}
	}
	if c.MaxBatch < 1 {
		return ValidationError{Field: "rollup.max_batch", Msg: "must be at least 1"}
	}
	if c.MaxBatch > 10_000 {
		return ValidationError{Field: "rollup.max_batch", Msg: "must not exceed 10000"}
	}
	return nil
}
