package config

import (
	"testing"
	"time"
)

func TestDefaultHipChatConfig_Values(t *testing.T) {
	cfg := DefaultHipChatConfig()

	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.BaseURL != "https://api.hipchat.com" {
		t.Errorf("unexpected BaseURL: %s", cfg.BaseURL)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("unexpected Timeout: %s", cfg.Timeout)
	}
	if cfg.Notify {
		t.Error("expected Notify to be false by default")
	}
	if cfg.MessageFmt == "" {
		t.Error("expected non-empty MessageFmt")
	}
}

func TestHipChatConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := DefaultHipChatConfig()
	cfg.Enabled = false
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error for disabled config, got: %v", err)
	}
}

func TestHipChatConfig_ValidateEnabledRequiresAuthToken(t *testing.T) {
	cfg := DefaultHipChatConfig()
	cfg.Enabled = true
	cfg.RoomID = "123456"
	// AuthToken intentionally empty
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for missing auth_token")
	}
	if !IsValidationError(err) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestHipChatConfig_ValidateEnabledRequiresRoomID(t *testing.T) {
	cfg := DefaultHipChatConfig()
	cfg.Enabled = true
	cfg.AuthToken = "tok-abc"
	// RoomID intentionally empty
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for missing room_id")
	}
	if !IsValidationError(err) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestHipChatConfig_ValidateTimeoutTooShort(t *testing.T) {
	cfg := DefaultHipChatConfig()
	cfg.Enabled = true
	cfg.AuthToken = "tok-abc"
	cfg.RoomID = "123456"
	cfg.Timeout = 500 * time.Millisecond
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for timeout < 1s")
	}
	if !IsValidationError(err) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestHipChatConfig_ValidateEnabledWithValidFields(t *testing.T) {
	cfg := DefaultHipChatConfig()
	cfg.Enabled = true
	cfg.AuthToken = "tok-abc"
	cfg.RoomID = "123456"
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error for valid config, got: %v", err)
	}
}
