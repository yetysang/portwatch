package config

import (
	"testing"
	"time"
)

func TestDefaultOpsGenieConfig_Values(t *testing.T) {
	cfg := DefaultOpsGenieConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.Priority != "P3" {
		t.Errorf("expected default priority P3, got %q", cfg.Priority)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", cfg.Timeout)
	}
	if cfg.APIKey != "" {
		t.Errorf("expected empty api_key by default, got %q", cfg.APIKey)
	}
}

func TestOpsGenieConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := DefaultOpsGenieConfig()
	cfg.Enabled = false
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestOpsGenieConfig_ValidateEnabledRequiresAPIKey(t *testing.T) {
	cfg := DefaultOpsGenieConfig()
	cfg.Enabled = true
	cfg.APIKey = ""
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error when api_key is empty")
	}
	if !IsValidationError(err) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestOpsGenieConfig_ValidateInvalidPriority(t *testing.T) {
	cfg := DefaultOpsGenieConfig()
	cfg.Enabled = true
	cfg.APIKey = "test-key"
	cfg.Priority = "HIGH"
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for invalid priority")
	}
	if !IsValidationError(err) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestOpsGenieConfig_ValidateEnabledWithValidFields(t *testing.T) {
	cfg := DefaultOpsGenieConfig()
	cfg.Enabled = true
	cfg.APIKey = "my-api-key"
	cfg.Priority = "P2"
	cfg.Timeout = 10 * time.Second
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for valid config, got %v", err)
	}
}

func TestOpsGenieConfig_ValidateZeroTimeoutInvalid(t *testing.T) {
	cfg := DefaultOpsGenieConfig()
	cfg.Enabled = true
	cfg.APIKey = "my-api-key"
	cfg.Timeout = 0
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for zero timeout")
	}
	if !IsValidationError(err) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}
