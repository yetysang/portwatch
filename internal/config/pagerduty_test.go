package config

import (
	"testing"
	"time"
)

func TestDefaultPagerDutyConfig_Values(t *testing.T) {
	cfg := DefaultPagerDutyConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.RoutingKey != "" {
		t.Errorf("expected empty RoutingKey, got %q", cfg.RoutingKey)
	}
	if cfg.Severity != "error" {
		t.Errorf("expected severity %q, got %q", "error", cfg.Severity)
	}
	if cfg.Timeout != 10*time.Second {
		t.Errorf("expected timeout 10s, got %v", cfg.Timeout)
	}
}

func TestPagerDutyConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := DefaultPagerDutyConfig()
	cfg.Enabled = false
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestPagerDutyConfig_ValidateEnabledRequiresRoutingKey(t *testing.T) {
	cfg := DefaultPagerDutyConfig()
	cfg.Enabled = true
	cfg.RoutingKey = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing routing_key")
	}
}

func TestPagerDutyConfig_ValidateEnabledWithRoutingKey(t *testing.T) {
	cfg := DefaultPagerDutyConfig()
	cfg.Enabled = true
	cfg.RoutingKey = "abc123"
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestPagerDutyConfig_ValidateInvalidSeverity(t *testing.T) {
	cfg := DefaultPagerDutyConfig()
	cfg.Enabled = true
	cfg.RoutingKey = "abc123"
	cfg.Severity = "unknown"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid severity")
	}
}

func TestPagerDutyConfig_ValidateAllSeverities(t *testing.T) {
	for _, sev := range []string{"critical", "error", "warning", "info"} {
		cfg := DefaultPagerDutyConfig()
		cfg.Enabled = true
		cfg.RoutingKey = "key"
		cfg.Severity = sev
		if err := cfg.Validate(); err != nil {
			t.Errorf("severity %q: unexpected error: %v", sev, err)
		}
	}
}

func TestPagerDutyConfig_ValidateZeroTimeout(t *testing.T) {
	cfg := DefaultPagerDutyConfig()
	cfg.Enabled = true
	cfg.RoutingKey = "key"
	cfg.Timeout = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero timeout")
	}
}
