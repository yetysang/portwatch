package config

import (
	"testing"
	"time"
)

func TestDefaultNewRelicConfig_Values(t *testing.T) {
	cfg := DefaultNewRelicConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false")
	}
	if cfg.EventType != "PortWatchEvent" {
		t.Errorf("unexpected EventType: %q", cfg.EventType)
	}
	if cfg.Region != "us" {
		t.Errorf("unexpected Region: %q", cfg.Region)
	}
	if cfg.Timeout != 10*time.Second {
		t.Errorf("unexpected Timeout: %v", cfg.Timeout)
	}
}

func TestNewRelicConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := DefaultNewRelicConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got: %v", err)
	}
}

func TestNewRelicConfig_ValidateEnabledRequiresAPIKey(t *testing.T) {
	cfg := DefaultNewRelicConfig()
	cfg.Enabled = true
	cfg.AccountID = "12345"
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for missing api_key")
	}
	if !IsValidationError(err) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestNewRelicConfig_ValidateEnabledRequiresAccountID(t *testing.T) {
	cfg := DefaultNewRelicConfig()
	cfg.Enabled = true
	cfg.APIKey = "abc123"
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for missing account_id")
	}
	if !IsValidationError(err) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestNewRelicConfig_ValidateInvalidRegion(t *testing.T) {
	cfg := DefaultNewRelicConfig()
	cfg.Enabled = true
	cfg.APIKey = "abc123"
	cfg.AccountID = "12345"
	cfg.Region = "ap"
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for invalid region")
	}
	if !IsValidationError(err) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestNewRelicConfig_ValidateEnabledWithValidFields(t *testing.T) {
	cfg := DefaultNewRelicConfig()
	cfg.Enabled = true
	cfg.APIKey = "abc123"
	cfg.AccountID = "12345"
	for _, region := range []string{"us", "eu"} {
		cfg.Region = region
		if err := cfg.Validate(); err != nil {
			t.Errorf("region %q: unexpected error: %v", region, err)
		}
	}
}
