package config

import "fmt"

// BearerTokenConfig holds configuration for outbound bearer token authentication.
type BearerTokenConfig struct {
	Enabled bool   `yaml:"enabled"`
	Token   string `yaml:"token"`
	Header  string `yaml:"header"`
}

// DefaultBearerTokenConfig returns a BearerTokenConfig with sensible defaults.
func DefaultBearerTokenConfig() BearerTokenConfig {
	return BearerTokenConfig{
		Enabled: false,
		Token:   "",
		Header:  "Authorization",
	}
}

// Validate checks that the BearerTokenConfig is consistent.
func (c BearerTokenConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Token == "" {
		return fmt.Errorf("bearer_token: token must not be empty when enabled")
	}
	if c.Header == "" {
		return fmt.Errorf("bearer_token: header must not be empty when enabled")
	}
	return nil
}
