package config

import (
	"testing"
	"time"
)

func TestDefaultJitterConfig_Values(t *testing.T) {
	cfg := DefaultJitterConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.MaxJitter != 500*time.Millisecond {
		t.Errorf("expected MaxJitter 500ms, got %s", cfg.MaxJitter)
	}
}

func TestJitterConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := JitterConfig{Enabled: false, MaxJitter: -1}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestJitterConfig_ValidateEnabledRequiresPositiveJitter(t *testing.T) {
	cfg := JitterConfig{Enabled: true, MaxJitter: 0}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero MaxJitter")
	}
}

func TestJitterConfig_ValidateNegativeJitter(t *testing.T) {
	cfg := JitterConfig{Enabled: true, MaxJitter: -1 * time.Second}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative MaxJitter")
	}
}

func TestJitterConfig_ValidateJitterTooLong(t *testing.T) {
	cfg := JitterConfig{Enabled: true, MaxJitter: 31 * time.Second}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for MaxJitter exceeding 30s")
	}
}

func TestJitterConfig_ValidateEnabledWithValidJitter(t *testing.T) {
	cfg := JitterConfig{Enabled: true, MaxJitter: 250 * time.Millisecond}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestJitterConfig_ValidateMaxBoundary(t *testing.T) {
	cfg := JitterConfig{Enabled: true, MaxJitter: 30 * time.Second}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error at boundary 30s, got %v", err)
	}
}
