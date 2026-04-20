package config

import (
	"fmt"
	"time"
)

// HealthCheckConfig controls the optional HTTP health-check endpoint.
type HealthCheckConfig struct {
	Enabled     bool          `toml:"enabled" json:"enabled"`
	ListenAddr  string        `toml:"listen_addr" json:"listen_addr"`
	Path        string        `toml:"path" json:"path"`
	ReadTimeout time.Duration `toml:"read_timeout" json:"read_timeout"`
}

// DefaultHealthCheckConfig returns a HealthCheckConfig with sensible defaults.
func DefaultHealthCheckConfig() HealthCheckConfig {
	return HealthCheckConfig{
		Enabled:     false,
		ListenAddr:  ":9090",
		Path:        "/healthz",
		ReadTimeout: 5 * time.Second,
	}
}

// Validate returns an error if the config is invalid.
func (c HealthCheckConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.ListenAddr == "" {
		return fmt.Errorf("healthcheck: listen_addr must not be empty")
	}
	if c.Path == "" || c.Path[0] != '/' {
		return fmt.Errorf("healthcheck: path must start with '/'")
	}
	if c.ReadTimeout <= 0 {
		return fmt.Errorf("healthcheck: read_timeout must be positive")
	}
	if c.ReadTimeout > 60*time.Second {
		return fmt.Errorf("healthcheck: read_timeout must not exceed 60s")
	}
	return nil
}
