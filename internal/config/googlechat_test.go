package config

import (
	"testing"
	"time"
)

func TestDefaultGoogleChatConfig_Values(t *testing.T) {
	cfg := DefaultGoogleChatConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false")
	}
	if cfg.WebhookURL != "" {
		t.Errorf("expected empty WebhookURL, got %q", cfg.WebhookURL)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("expected 5s timeout, got %v", cfg.Timeout)
	}
}

func TestGoogleChatConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := DefaultGoogleChatConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestGoogleChatConfig_ValidateEnabledRequiresURL(t *testing.T) {
	cfg := DefaultGoogleChatConfig()
	cfg.Enabled = true
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when webhook_url is empty")
	}
}

func TestGoogleChatConfig_ValidateEnabledWithURL(t *testing.T) {
	cfg := DefaultGoogleChatConfig()
	cfg.Enabled = true
	cfg.WebhookURL = "https://chat.googleapis.com/v1/spaces/xxx/messages?key=yyy"
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGoogleChatConfig_ValidateTimeoutTooShort(t *testing.T) {
	cfg := DefaultGoogleChatConfig()
	cfg.Enabled = true
	cfg.WebhookURL = "https://chat.googleapis.com/v1/spaces/xxx/messages?key=yyy"
	cfg.Timeout = 500 * time.Millisecond
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for timeout < 1s")
	}
}

func TestGoogleChatConfig_ValidateTimeoutTooLong(t *testing.T) {
	cfg := DefaultGoogleChatConfig()
	cfg.Enabled = true
	cfg.WebhookURL = "https://chat.googleapis.com/v1/spaces/xxx/messages?key=yyy"
	cfg.Timeout = 60 * time.Second
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for timeout > 30s")
	}
}
