package config

import "testing"

func TestDefaultTelegramConfig_Values(t *testing.T) {
	c := DefaultTelegramConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if c.BotToken != "" {
		t.Errorf("expected empty BotToken, got %q", c.BotToken)
	}
	if c.ChatID != "" {
		t.Errorf("expected empty ChatID, got %q", c.ChatID)
	}
	if c.ParseMode != "Markdown" {
		t.Errorf("expected ParseMode \"Markdown\", got %q", c.ParseMode)
	}
}

func TestTelegramConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	c := TelegramConfig{Enabled: false}
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestTelegramConfig_ValidateEnabledRequiresBotToken(t *testing.T) {
	c := TelegramConfig{Enabled: true, ChatID: "123456"}
	err := c.Validate()
	if err == nil {
		t.Fatal("expected error for missing bot_token")
	}
	if !IsValidationError(err) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestTelegramConfig_ValidateEnabledRequiresChatID(t *testing.T) {
	c := TelegramConfig{Enabled: true, BotToken: "abc:def"}
	err := c.Validate()
	if err == nil {
		t.Fatal("expected error for missing chat_id")
	}
	if !IsValidationError(err) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestTelegramConfig_ValidateInvalidParseMode(t *testing.T) {
	c := TelegramConfig{Enabled: true, BotToken: "abc:def", ChatID: "123", ParseMode: "INVALID"}
	err := c.Validate()
	if err == nil {
		t.Fatal("expected error for invalid parse_mode")
	}
	if !IsValidationError(err) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestTelegramConfig_ValidateValidConfig(t *testing.T) {
	c := TelegramConfig{Enabled: true, BotToken: "abc:def", ChatID: "123456", ParseMode: "HTML"}
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error for valid config, got %v", err)
	}
}

func TestTelegramConfig_ValidateEmptyParseMode(t *testing.T) {
	c := TelegramConfig{Enabled: true, BotToken: "abc:def", ChatID: "123456", ParseMode: ""}
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error for empty parse_mode, got %v", err)
	}
}
