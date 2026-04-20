package config

import (
	"testing"
	"time"
)

func TestDefaultHealthCheckConfig_Values(t *testing.T) {
	cfg := DefaultHealthCheckConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.ListenAddr != ":9090" {
		t.Errorf("unexpected ListenAddr: %s", cfg.ListenAddr)
	}
	if cfg.Path != "/healthz" {
		t.Errorf("unexpected Path: %s", cfg.Path)
	}
	if cfg.ReadTimeout != 5*time.Second {
		t.Errorf("unexpected ReadTimeout: %v", cfg.ReadTimeout)
	}
}

func TestHealthCheckConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := HealthCheckConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got: %v", err)
	}
}

func TestHealthCheckConfig_ValidateEnabledRequiresListenAddr(t *testing.T) {
	cfg := DefaultHealthCheckConfig()
	cfg.Enabled = true
	cfg.ListenAddr = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty listen_addr")
	}
}

func TestHealthCheckConfig_ValidatePathMustStartWithSlash(t *testing.T) {
	cfg := DefaultHealthCheckConfig()
	cfg.Enabled = true
	cfg.Path = "healthz"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for path without leading slash")
	}
}

func TestHealthCheckConfig_ValidateReadTimeoutTooLong(t *testing.T) {
	cfg := DefaultHealthCheckConfig()
	cfg.Enabled = true
	cfg.ReadTimeout = 120 * time.Second
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for read_timeout > 60s")
	}
}

func TestHealthCheckConfig_ValidateValidConfig(t *testing.T) {
	cfg := DefaultHealthCheckConfig()
	cfg.Enabled = true
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for valid config, got: %v", err)
	}
}
