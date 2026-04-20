package config

import "testing"

func TestDefaultSignalConfig_Values(t *testing.T) {
	cfg := DefaultSignalConfig()
	if cfg.GracePeriodSeconds != 5 {
		t.Errorf("expected GracePeriodSeconds=5, got %d", cfg.GracePeriodSeconds)
	}
	if !cfg.ReloadOnHUP {
		t.Error("expected ReloadOnHUP=true")
	}
}

func TestSignalConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := DefaultSignalConfig()
	cfg.GracePeriodSeconds = 0
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error for zero grace period: %v", err)
	}
}

func TestSignalConfig_ValidateNegativeGracePeriod(t *testing.T) {
	cfg := DefaultSignalConfig()
	cfg.GracePeriodSeconds = -1
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for negative grace period")
	}
	if !IsValidationError(err) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestSignalConfig_ValidateGracePeriodTooLong(t *testing.T) {
	cfg := DefaultSignalConfig()
	cfg.GracePeriodSeconds = 121
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for grace period > 120")
	}
	if !IsValidationError(err) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestSignalConfig_ValidateMaxGracePeriod(t *testing.T) {
	cfg := DefaultSignalConfig()
	cfg.GracePeriodSeconds = 120
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error for max grace period: %v", err)
	}
}

func TestSignalConfig_ValidateReloadOnHUPFalse(t *testing.T) {
	cfg := DefaultSignalConfig()
	cfg.ReloadOnHUP = false
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error when ReloadOnHUP=false: %v", err)
	}
}
