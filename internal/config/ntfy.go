package config

import "fmt"

// NtfyConfig holds ntfy.sh alert integration settings.
type NtfyConfig struct {
	Enabled   bool   `yaml:"enabled"`
	ServerURL string `yaml:"server_url"`
	Topic     string `yaml:"topic"`
	Token     string `yaml:"token,omitempty"`
	Priority  string `yaml:"priority"`
}

// DefaultNtfyConfig returns a NtfyConfig with sensible defaults.
func DefaultNtfyConfig() NtfyConfig {
	return NtfyConfig{
		Enabled:   false,
		ServerURL: "https://ntfy.sh",
		Topic:     "portwatch",
		Priority:  "default",
	}
}

// Validate checks that the NtfyConfig is valid when enabled.
func (c NtfyConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.ServerURL == "" {
		return fmt.Errorf("ntfy: server_url is required when enabled")
	}
	if c.Topic == "" {
		return fmt.Errorf("ntfy: topic is required when enabled")
	}
	allowed := map[string]bool{
		"max":     true,
		"urgent":  true,
		"high":    true,
		"default": true,
		"low":     true,
		"min":     true,
	}
	if !allowed[c.Priority] {
		return fmt.Errorf("ntfy: invalid priority %q", c.Priority)
	}
	return nil
}
