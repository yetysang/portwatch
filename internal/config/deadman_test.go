package config

import (
	"testing"
	"time"
)

func TestDefaultDeadManConfig_Values(t *testing.T) {
	cfg := DefaultDeadManConfig()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Interval != 60*time.Second {
		t.Errorf("expected Interval=60s, got %s", cfg.Interval)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("expected Timeout=5s, got %s", cfg.Timeout)
	}
}

func TestDeadManConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := DeadManConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestDeadManConfig_ValidateEnabledRequiresURL(t *testing.T) {
	cfg := DeadManConfig{
		Enabled:  true,
		Interval: 30 * time.Second,
		Timeout:  5 * time.Second,
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing URL")
	}
}

func TestDeadManConfig_ValidateIntervalTooShort(t *testing.T) {
	cfg := DeadManConfig{
		Enabled:  true,
		URL:      "https://hc-ping.example.com/abc",
		Interval: 5 * time.Second,
		Timeout:  2 * time.Second,
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for interval < 10s")
	}
}

func TestDeadManConfig_ValidateIntervalTooLong(t *testing.T) {
	cfg := DeadManConfig{
		Enabled:  true,
		URL:      "https://hc-ping.example.com/abc",
		Interval: 25 * time.Hour,
		Timeout:  5 * time.Second,
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for interval > 24h")
	}
}

func TestDeadManConfig_ValidateTimeoutMustBePositive(t *testing.T) {
	cfg := DeadManConfig{
		Enabled:  true,
		URL:      "https://hc-ping.example.com/abc",
		Interval: 60 * time.Second,
		Timeout:  0,
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero timeout")
	}
}

func TestDeadManConfig_ValidateEnabledWithValidFields(t *testing.T) {
	cfg := DeadManConfig{
		Enabled:  true,
		URL:      "https://hc-ping.example.com/abc123",
		Interval: 60 * time.Second,
		Timeout:  5 * time.Second,
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for valid config, got %v", err)
	}
}
