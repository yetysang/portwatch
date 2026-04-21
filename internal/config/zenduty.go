package config

import "fmt"

// ZendutyConfig holds settings for the Zenduty alert handler.
type ZendutyConfig struct {
	Enabled    bool   `toml:"enabled"`
	APIKey     string `toml:"api_key"`
	ServiceID  string `toml:"service_id"`
	IntegrationID string `toml:"integration_id"`
	AlertType  string `toml:"alert_type"` // "info", "warning", "critical"
}

// DefaultZendutyConfig returns a ZendutyConfig with safe defaults.
func DefaultZendutyConfig() ZendutyConfig {
	return ZendutyConfig{
		Enabled:   false,
		AlertType: "warning",
	}
}

// Validate returns an error if the ZendutyConfig is invalid.
func (c ZendutyConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.APIKey == "" {
		return fmt.Errorf("zenduty: api_key is required when enabled")
	}
	if c.ServiceID == "" {
		return fmt.Errorf("zenduty: service_id is required when enabled")
	}
	if c.IntegrationID == "" {
		return fmt.Errorf("zenduty: integration_id is required when enabled")
	}
	switch c.AlertType {
	case "info", "warning", "critical":
		// valid
	default:
		return fmt.Errorf("zenduty: alert_type must be one of info, warning, critical; got %q", c.AlertType)
	}
	return nil
}
