package config

import "fmt"

// MatrixConfig holds configuration for the Matrix alert handler.
type MatrixConfig struct {
	Enabled    bool   `toml:"enabled" yaml:"enabled"`
	Homeserver string `toml:"homeserver" yaml:"homeserver"`
	Token      string `toml:"token" yaml:"token"`
	RoomID     string `toml:"room_id" yaml:"room_id"`
}

// DefaultMatrixConfig returns a MatrixConfig with safe defaults.
func DefaultMatrixConfig() MatrixConfig {
	return MatrixConfig{
		Enabled:    false,
		Homeserver: "",
		Token:      "",
		RoomID:     "",
	}
}

// Validate checks that the Matrix configuration is self-consistent.
func (c MatrixConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Homeserver == "" {
		return fmt.Errorf("matrix: homeserver must not be empty when enabled")
	}
	if c.Token == "" {
		return fmt.Errorf("matrix: token must not be empty when enabled")
	}
	if c.RoomID == "" {
		return fmt.Errorf("matrix: room_id must not be empty when enabled")
	}
	return nil
}
