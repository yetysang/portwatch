package config

import (
	"fmt"
	"time"
)

// HipChatConfig holds configuration for the HipChat alert handler.
type HipChatConfig struct {
	Enabled   bool          `toml:"enabled"`
	AuthToken string        `toml:"auth_token"`
	RoomID    string        `toml:"room_id"`
	BaseURL   string        `toml:"base_url"`
	Color     string        `toml:"color"`
	Notify    bool          `toml:"notify"`
	Timeout   time.Duration `toml:"timeout"`
}

// DefaultHipChatConfig returns a HipChatConfig with sensible defaults.
func DefaultHipChatConfig() HipChatConfig {
	return HipChatConfig{
		Enabled:  false,
		BaseURL:  "https://api.hipchat.com",
		Color:    "yellow",
		Notify:   false,
		Timeout:  10 * time.Second,
	}
}

// Validate checks that the HipChatConfig is consistent.
func (c HipChatConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.AuthToken == "" {
		return ValidationError{Field: "auth_token", Msg: "auth_token is required when hipchat is enabled"}
	}
	if c.RoomID == "" {
		return ValidationError{Field: "room_id", Msg: "room_id is required when hipchat is enabled"}
	}
	if c.Timeout < time.Second {
		return ValidationError{Field: "timeout", Msg: fmt.Sprintf("timeout must be at least 1s, got %s", c.Timeout)}
	}
	return nil
}
