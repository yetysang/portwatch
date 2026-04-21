package config

import (
	"fmt"
	"time"
)

// ScriptConfig holds configuration for the script alert handler.
type ScriptConfig struct {
	Enabled  bool          `toml:"enabled"`
	Path     string        `toml:"path"`
	Args     []string      `toml:"args"`
	Timeout  time.Duration `toml:"timeout"`
	EnvVars  []string      `toml:"env_vars"`
}

// DefaultScriptConfig returns a ScriptConfig with sensible defaults.
func DefaultScriptConfig() ScriptConfig {
	return ScriptConfig{
		Enabled: false,
		Path:    "",
		Args:    []string{},
		Timeout: 10 * time.Second,
		EnvVars: []string{},
	}
}

// Validate checks that the ScriptConfig is valid.
func (c ScriptConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Path == "" {
		return fmt.Errorf("script.path must be set when script alerts are enabled")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("script.timeout must be positive, got %s", c.Timeout)
	}
	if c.Timeout > 5*time.Minute {
		return fmt.Errorf("script.timeout must not exceed 5m, got %s", c.Timeout)
	}
	return nil
}
