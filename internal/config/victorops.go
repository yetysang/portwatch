package config

import (
	"fmt"
	"time"
)

// VictorOpsConfig holds configuration for the VictorOps (Splunk On-Call) alert handler.
type VictorOpsConfig struct {
	Enabled    bool          `yaml:"enabled"`
	URL        string        `yaml:"url"`
	RoutingKey string        `yaml:"routing_key"`
	Timeout    time.Duration `yaml:"timeout"`
}

// DefaultVictorOpsConfig returns a VictorOpsConfig with sensible defaults.
func DefaultVictorOpsConfig() VictorOpsConfig {
	return VictorOpsConfig{
		Enabled:    false,
		URL:        "https://alert.victorops.com/integrations/generic/20131114/alert",
		RoutingKey: "",
		Timeout:    10 * time.Second,
	}
}

// Validate checks the VictorOpsConfig for correctness.
func (c VictorOpsConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.URL == "" {
		return fmt.Errorf("victorops: url is required when enabled")
	}
	if c.RoutingKey == "" {
		return fmt.Errorf("victorops: routing_key is required when enabled")
	}
	if c.Timeout < time.Second {
		return fmt.Errorf("victorops: timeout must be at least 1s")
	}
	return nil
}
