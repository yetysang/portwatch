package config

import (
	"fmt"
	"time"
)

// HipChatConfig holds configuration for the HipChat alert handler.
type HipChatConfig struct {
	Enabled    bool          `toml:"enabled" yaml:"enabled"`
	AuthToken  string        `toml:"auth_token" yaml:"auth_token"`
	RoomID     string        `toml:"room_id" yaml:"room_id"`
	BaseURL    string        `toml:"base_url" yaml:"base_url"`
	Timeout    time.Duration `toml:"timeout" yaml:"timeout"`
	Notify     bool          `toml:"notify" yaml:"notify"`
	MessageFmt string        `toml:"message_fmt" yaml:"message_fmt"`
}

// DefaultHipChatConfig returns a HipChatConfig with sensible defaults.
func DefaultHipChatConfig() HipChatConfig {
	return HipChatConfig{
		Enabled:    false,
		BaseURL:    "https://api.hipchat.com",
		Timeout:    5 * time.Second,
		Notify:     false,
		MessageFmt: "[portwatch] {action} {proto}:{port} on {host}",
	}
}

// Validate checks that the HipChatConfig is internally consistent.
func (c HipChatConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.AuthToken == "" {
		return ValidationError{Field: "auth_token", Msg: "required when hipchat is enabled"}
	}
	if c.RoomID == "" {
		return ValidationError{Field: "room_id", Msg: "required when hipchat is enabled"}
	}
	if c.Timeout < time.Second {
		return ValidationError{Field: "timeout", Msg: fmt.Sprintf("must be at least 1s, got %s", c.Timeout)}
	}
	return nil
}
