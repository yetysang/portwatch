package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func TestDefaultGotifyConfig_Values(t *testing.T) {
	cfg := config.DefaultGotifyConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.URL != "" {
		t.Errorf("expected empty URL, got %q", cfg.URL)
	}
	if cfg.Token != "" {
		t.Errorf("expected empty Token, got %q", cfg.Token)
	}
	if cfg.Priority != 5 {
		t.Errorf("expected default Priority 5, got %d", cfg.Priority)
	}
}

func TestGotifyConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := config.DefaultGotifyConfig()
	cfg.Enabled = false
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestGotifyConfig_ValidateEnabledRequiresURL(t *testing.T) {
	cfg := config.DefaultGotifyConfig()
	cfg.Enabled = true
	cfg.Token = "tok"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when URL is empty")
	}
}

func TestGotifyConfig_ValidateEnabledRequiresToken(t *testing.T) {
	cfg := config.DefaultGotifyConfig()
	cfg.Enabled = true
	cfg.URL = "http://gotify.example.com"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when Token is empty")
	}
}

func TestGotifyConfig_ValidateEnabledWithBothFields(t *testing.T) {
	cfg := config.DefaultGotifyConfig()
	cfg.Enabled = true
	cfg.URL = "http://gotify.example.com"
	cfg.Token = "mytoken"
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error with valid config, got %v", err)
	}
}

func TestGotifyConfig_ValidateInvalidPriority(t *testing.T) {
	cfg := config.DefaultGotifyConfig()
	cfg.Enabled = true
	cfg.URL = "http://gotify.example.com"
	cfg.Token = "mytoken"
	cfg.Priority = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for priority < 1")
	}
}
