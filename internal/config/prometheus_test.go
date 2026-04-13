package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func TestDefaultPrometheusConfig_Values(t *testing.T) {
	cfg := config.DefaultPrometheusConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.ListenAddr != ":9090" {
		t.Errorf("expected ListenAddr ':9090', got %q", cfg.ListenAddr)
	}
	if cfg.Path != "/metrics" {
		t.Errorf("expected Path '/metrics', got %q", cfg.Path)
	}
}

func TestPrometheusConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := config.DefaultPrometheusConfig()
	cfg.Enabled = false
	cfg.ListenAddr = ""
	cfg.Path = ""
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestPrometheusConfig_ValidateEnabledRequiresListenAddr(t *testing.T) {
	cfg := config.DefaultPrometheusConfig()
	cfg.Enabled = true
	cfg.ListenAddr = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when ListenAddr is empty")
	}
}

func TestPrometheusConfig_ValidateEnabledRequiresPath(t *testing.T) {
	cfg := config.DefaultPrometheusConfig()
	cfg.Enabled = true
	cfg.Path = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when Path is empty")
	}
}

func TestPrometheusConfig_ValidatePathMustStartWithSlash(t *testing.T) {
	cfg := config.DefaultPrometheusConfig()
	cfg.Enabled = true
	cfg.Path = "metrics"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when Path does not start with '/'")
	}
}

func TestPrometheusConfig_ValidateValidConfig(t *testing.T) {
	cfg := config.DefaultPrometheusConfig()
	cfg.Enabled = true
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for valid config, got %v", err)
	}
}
