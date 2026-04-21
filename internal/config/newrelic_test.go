package config

import (
	"testing"
	"time"
)

func TestDefaultNewRelicConfig_Values(t *testing.T) {
	cfg := DefaultNewRelicConfig()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Region != "US" {
		t.Errorf("expected Region=US, got %s", cfg.Region)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("expected Timeout=5s, got %v", cfg.Timeout)
	}
}

func TestNewRelicConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := DefaultNewRelicConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestNewRelicConfig_ValidateEnabledRequiresAPIKey(t *testing.T) {
	cfg := DefaultNewRelicConfig()
	cfg.Enabled = true
	cfg.AccountID = "123"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when api_key is missing")
	}
}

func TestNewRelicConfig_ValidateEnabledRequiresAccountID(t *testing.T) {
	cfg := DefaultNewRelicConfig()
	cfg.Enabled = true
	cfg.APIKey = "key"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when account_id is missing")
	}
}

func TestNewRelicConfig_ValidateInvalidRegion(t *testing.T) {
	cfg := DefaultNewRelicConfig()
	cfg.Enabled = true
	cfg.APIKey = "key"
	cfg.AccountID = "123"
	cfg.Region = "AP"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid region")
	}
}

func TestNewRelicConfig_ValidateEnabledWithValidFields(t *testing.T) {
	cfg := DefaultNewRelicConfig()
	cfg.Enabled = true
	cfg.APIKey = "key"
	cfg.AccountID = "123456"
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNewRelicConfig_ValidateEURegion(t *testing.T) {
	cfg := DefaultNewRelicConfig()
	cfg.Enabled = true
	cfg.APIKey = "key"
	cfg.AccountID = "123456"
	cfg.Region = "EU"
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error for EU region: %v", err)
	}
}
