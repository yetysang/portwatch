package config_test

import (
	"testing"

	"github.com/wander/portwatch/internal/config"
)

func TestDefaultNatsConfig_Values(t *testing.T) {
	cfg := config.DefaultNatsConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false")
	}
	if cfg.URL != "nats://localhost:4222" {
		t.Errorf("unexpected URL: %s", cfg.URL)
	}
	if cfg.Subject != "portwatch.events" {
		t.Errorf("unexpected Subject: %s", cfg.Subject)
	}
}

func TestNatsConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := config.NatsConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestNatsConfig_ValidateEnabledRequiresURL(t *testing.T) {
	cfg := config.NatsConfig{Enabled: true, Subject: "portwatch.events"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing URL")
	}
}

func TestNatsConfig_ValidateEnabledRequiresSubject(t *testing.T) {
	cfg := config.NatsConfig{Enabled: true, URL: "nats://localhost:4222"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing Subject")
	}
}

func TestNatsConfig_ValidateEnabledWithBothFields(t *testing.T) {
	cfg := config.NatsConfig{
		Enabled: true,
		URL:     "nats://localhost:4222",
		Subject: "portwatch.events",
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestNatsConfig_ValidateErrorMessage(t *testing.T) {
	cfg := config.NatsConfig{Enabled: true, URL: "nats://localhost:4222"}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error")
	}
	got := err.Error()
	if got != "nats: subject is required when enabled" {
		t.Errorf("unexpected error message: %s", got)
	}
}
