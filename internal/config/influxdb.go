package config

import (
	"fmt"
	"time"
)

// InfluxDBConfig holds configuration for the InfluxDB alert handler.
type InfluxDBConfig struct {
	Enabled      bool          `toml:"enabled" yaml:"enabled"`
	URL          string        `toml:"url" yaml:"url"`
	Token        string        `toml:"token" yaml:"token"`
	Org          string        `toml:"org" yaml:"org"`
	Bucket       string        `toml:"bucket" yaml:"bucket"`
	Measurement  string        `toml:"measurement" yaml:"measurement"`
	Timeout      time.Duration `toml:"timeout" yaml:"timeout"`
}

// DefaultInfluxDBConfig returns an InfluxDBConfig with sensible defaults.
func DefaultInfluxDBConfig() InfluxDBConfig {
	return InfluxDBConfig{
		Enabled:     false,
		Measurement: "portwatch_events",
		Timeout:     5 * time.Second,
	}
}

// Validate checks that the InfluxDB configuration is valid.
func (c InfluxDBConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return fmt.Errorf("influxdb: url is required when enabled")
	}
	if c.Token == "" {
		return fmt.Errorf("influxdb: token is required when enabled")
	}
	if c.Org == "" {
		return fmt.Errorf("influxdb: org is required when enabled")
	}
	if c.Bucket == "" {
		return fmt.Errorf("influxdb: bucket is required when enabled")
	}
	if c.Measurement == "" {
		return fmt.Errorf("influxdb: measurement must not be empty")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("influxdb: timeout must be positive")
	}
	return nil
}
