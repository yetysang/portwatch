package config

import (
	"testing"
	"time"
)

func TestDefaultGRPCConfig_Values(t *testing.T) {
	cfg := DefaultGRPCConfig()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Target != "localhost:50051" {
		t.Errorf("unexpected default target: %s", cfg.Target)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("unexpected default timeout: %s", cfg.Timeout)
	}
}

func TestGRPCConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := GRPCConfig{Enabled: false, Target: "", Timeout: 0}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error when disabled, got: %v", err)
	}
}

func TestGRPCConfig_ValidateEnabledRequiresTarget(t *testing.T) {
	cfg := GRPCConfig{Enabled: true, Target: "", Timeout: 5 * time.Second}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty target")
	}
}

func TestGRPCConfig_ValidateEnabledRequiresPositiveTimeout(t *testing.T) {
	cfg := GRPCConfig{Enabled: true, Target: "host:1234", Timeout: 0}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero timeout")
	}
}

func TestGRPCConfig_ValidateEnabledWithValidFields(t *testing.T) {
	cfg := GRPCConfig{Enabled: true, Target: "host:1234", Timeout: 3 * time.Second}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for valid config, got: %v", err)
	}
}

func TestGRPCConfig_ValidateNegativeTimeout(t *testing.T) {
	cfg := GRPCConfig{Enabled: true, Target: "host:9090", Timeout: -1 * time.Second}
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for negative timeout")
	}
}
