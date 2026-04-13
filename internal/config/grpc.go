package config

import (
	"fmt"
	"time"
)

// GRPCConfig holds configuration for the gRPC alert handler.
type GRPCConfig struct {
	Enabled bool          `toml:"enabled" yaml:"enabled"`
	Target  string        `toml:"target"  yaml:"target"`
	Timeout time.Duration `toml:"timeout" yaml:"timeout"`
}

// DefaultGRPCConfig returns a GRPCConfig with sensible defaults.
func DefaultGRPCConfig() GRPCConfig {
	return GRPCConfig{
		Enabled: false,
		Target:  "localhost:50051",
		Timeout: 5 * time.Second,
	}
}

// Validate returns an error if the configuration is invalid.
func (c GRPCConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Target == "" {
		return fmt.Errorf("grpc: target must not be empty when enabled")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("grpc: timeout must be positive, got %s", c.Timeout)
	}
	return nil
}
