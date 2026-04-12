package config

import "fmt"

// TelegramConfig holds configuration for the Telegram alert handler.
type TelegramConfig struct {
	Enabled  bool   `toml:"enabled" json:"enabled"`
	BotToken string `toml:"bot_token" json:"bot_token"`
	ChatID   string `toml:"chat_id" json:"chat_id"`
	ParseMode string `toml:"parse_mode" json:"parse_mode"`
}

// DefaultTelegramConfig returns a TelegramConfig with sensible defaults.
func DefaultTelegramConfig() TelegramConfig {
	return TelegramConfig{
		Enabled:   false,
		BotToken:  "",
		ChatID:    "",
		ParseMode: "Markdown",
	}
}

// Validate checks that the TelegramConfig is consistent.
// Returns a ValidationError if the config is invalid.
func (c TelegramConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.BotToken == "" {
		return &ValidationError{Field: "telegram.bot_token", Message: "bot_token is required when telegram alerts are enabled"}
	}
	if c.ChatID == "" {
		return &ValidationError{Field: "telegram.chat_id", Message: "chat_id is required when telegram alerts are enabled"}
	}
	validModes := map[string]bool{"Markdown": true, "MarkdownV2": true, "HTML": true, "": true}
	if !validModes[c.ParseMode] {
		return &ValidationError{
			Field:   "telegram.parse_mode",
			Message: fmt.Sprintf("parse_mode %q is not valid; expected Markdown, MarkdownV2, HTML, or empty", c.ParseMode),
		}
	}
	return nil
}
