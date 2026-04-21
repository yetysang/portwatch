package config

import (
	"testing"
	"time"
)

func TestDefaultRateLimitConfig_Values(t *testing.T) {
	cfg := DefaultRateLimitConfig()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
	if cfg.Cooldown != 30*time.Second {
		t.Errorf("expected Cooldown 30s, got %s", cfg.Cooldown)
	}
	if cfg.MaxBurst != 1 {
		t.Errorf("expected MaxBurst 1, got %d", cfg.MaxBurst)
	}
}

func TestRateLimitConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := RateLimitConfig{Enabled: false, Cooldown: 0, MaxBurst: 0}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestRateLimitConfig_ValidateCooldownTooShort(t *testing.T) {
	cfg := DefaultRateLimitConfig()
	cfg.Cooldown = 500 * time.Millisecond
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for cooldown < 1s")
	}
}

func TestRateLimitConfig_ValidateCooldownTooLong(t *testing.T) {
	cfg := DefaultRateLimitConfig()
	cfg.Cooldown = 25 * time.Hour
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for cooldown > 24h")
	}
}

func TestRateLimitConfig_ValidateMaxBurstTooLow(t *testing.T) {
	cfg := DefaultRateLimitConfig()
	cfg.MaxBurst = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for max_burst < 1")
	}
}

func TestRateLimitConfig_ValidateMaxBurstTooHigh(t *testing.T) {
	cfg := DefaultRateLimitConfig()
	cfg.MaxBurst = 101
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for max_burst > 100")
	}
}

func TestRateLimitConfig_ValidateValidConfig(t *testing.T) {
	cfg := RateLimitConfig{
		Enabled:  true,
		Cooldown: 10 * time.Second,
		MaxBurst: 3,
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for valid config, got %v", err)
	}
}
