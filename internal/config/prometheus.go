package config

import "fmt"

// PrometheusConfig holds settings for the Prometheus metrics handler.
type PrometheusConfig struct {
	// Enabled controls whether the Prometheus handler is active.
	Enabled bool `toml:"enabled" yaml:"enabled"`

	// ListenAddr is the address the metrics HTTP server listens on.
	ListenAddr string `toml:"listen_addr" yaml:"listen_addr"`

	// Path is the URL path that exposes metrics.
	Path string `toml:"path" yaml:"path"`
}

// DefaultPrometheusConfig returns a PrometheusConfig with sensible defaults.
func DefaultPrometheusConfig() PrometheusConfig {
	return PrometheusConfig{
		Enabled:    false,
		ListenAddr: ":9090",
		Path:       "/metrics",
	}
}

// Validate checks that the PrometheusConfig is internally consistent.
func (c PrometheusConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.ListenAddr == "" {
		return fmt.Errorf("prometheus: listen_addr must not be empty when enabled")
	}
	if c.Path == "" {
		return fmt.Errorf("prometheus: path must not be empty when enabled")
	}
	if c.Path[0] != '/' {
		return fmt.Errorf("prometheus: path must begin with '/', got %q", c.Path)
	}
	return nil
}
