package config

import (
	"fmt"
	"time"
)

// PagerDutyConfig holds configuration for the PagerDuty alert handler.
type PagerDutyConfig struct {
	Enabled    bool          `yaml:"enabled"`
	RoutingKey string        `yaml:"routing_key"`
	Severity   string        `yaml:"severity"`
	Timeout    time.Duration `yaml:"timeout"`
}

// DefaultPagerDutyConfig returns a PagerDutyConfig with sensible defaults.
func DefaultPagerDutyConfig() PagerDutyConfig {
	return PagerDutyConfig{
		Enabled:    false,
		RoutingKey: "",
		Severity:   "error",
		Timeout:    10 * time.Second,
	}
}

// Validate checks that the PagerDutyConfig is valid.
func (c PagerDutyConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.RoutingKey == "" {
		return fmt.Errorf("pagerduty: routing_key is required when enabled")
	}
	valid := map[string]bool{
		"critical": true,
		"error":    true,
		"warning":  true,
		"info":     true,
	}
	if !valid[c.Severity] {
		return fmt.Errorf("pagerduty: severity must be one of critical, error, warning, info; got %q", c.Severity)
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("pagerduty: timeout must be positive")
	}
	return nil
}
