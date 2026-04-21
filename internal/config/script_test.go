package config

import (
	"testing"
	"time"
)

func TestDefaultScriptConfig_Values(t *testing.T) {
	cfg := DefaultScriptConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.Path != "" {
		t.Errorf("expected empty Path, got %q", cfg.Path)
	}
	if cfg.Timeout != 10*time.Second {
		t.Errorf("expected 10s timeout, got %s", cfg.Timeout)
	}
}

func TestScriptConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := DefaultScriptConfig()
	cfg.Enabled = false
	cfg.Path = ""
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestScriptConfig_ValidateEnabledRequiresPath(t *testing.T) {
	cfg := DefaultScriptConfig()
	cfg.Enabled = true
	cfg.Path = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when path is empty")
	}
}

func TestScriptConfig_ValidateEnabledWithPath(t *testing.T) {
	cfg := DefaultScriptConfig()
	cfg.Enabled = true
	cfg.Path = "/usr/local/bin/notify.sh"
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestScriptConfig_ValidateTimeoutTooLong(t *testing.T) {
	cfg := DefaultScriptConfig()
	cfg.Enabled = true
	cfg.Path = "/usr/local/bin/notify.sh"
	cfg.Timeout = 10 * time.Minute
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when timeout exceeds 5m")
	}
}

func TestScriptConfig_ValidateZeroTimeoutInvalid(t *testing.T) {
	cfg := DefaultScriptConfig()
	cfg.Enabled = true
	cfg.Path = "/usr/local/bin/notify.sh"
	cfg.Timeout = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when timeout is zero")
	}
}
