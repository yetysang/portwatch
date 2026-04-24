package config_test

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/config"
)

func TestDefaultVictorOpsConfig_Values(t *testing.T) {
	cfg := config.DefaultVictorOpsConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false")
	}
	if cfg.URL == "" {
		t.Error("expected non-empty default URL")
	}
	if cfg.RoutingKey != "" {
		t.Errorf("expected empty RoutingKey, got %q", cfg.RoutingKey)
	}
	if cfg.Timeout != 10*time.Second {
		t.Errorf("expected 10s timeout, got %v", cfg.Timeout)
	}
}

func TestVictorOpsConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := config.VictorOpsConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestVictorOpsConfig_ValidateEnabledRequiresURL(t *testing.T) {
	cfg := config.VictorOpsConfig{
		Enabled:    true,
		URL:        "",
		RoutingKey: "team-key",
		Timeout:    5 * time.Second,
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing URL")
	}
}

func TestVictorOpsConfig_ValidateEnabledRequiresRoutingKey(t *testing.T) {
	cfg := config.VictorOpsConfig{
		Enabled:    true,
		URL:        "https://alert.victorops.com/integrations/generic/20131114/alert",
		RoutingKey: "",
		Timeout:    5 * time.Second,
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing routing_key")
	}
}

func TestVictorOpsConfig_ValidateTimeoutTooShort(t *testing.T) {
	cfg := config.VictorOpsConfig{
		Enabled:    true,
		URL:        "https://alert.victorops.com/integrations/generic/20131114/alert",
		RoutingKey: "team-key",
		Timeout:    500 * time.Millisecond,
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for timeout below 1s")
	}
}

func TestVictorOpsConfig_ValidateEnabledWithValidFields(t *testing.T) {
	cfg := config.VictorOpsConfig{
		Enabled:    true,
		URL:        "https://alert.victorops.com/integrations/generic/20131114/alert",
		RoutingKey: "team-key",
		Timeout:    10 * time.Second,
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
