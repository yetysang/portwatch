package config

import "testing"

func TestDefaultNtfyConfig_Values(t *testing.T) {
	cfg := DefaultNtfyConfig()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.ServerURL != "https://ntfy.sh" {
		t.Errorf("ServerURL = %q, want %q", cfg.ServerURL, "https://ntfy.sh")
	}
	if cfg.Topic != "portwatch" {
		t.Errorf("Topic = %q, want %q", cfg.Topic, "portwatch")
	}
	if cfg.Priority != "default" {
		t.Errorf("Priority = %q, want %q", cfg.Priority, "default")
	}
}

func TestNtfyConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := NtfyConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestNtfyConfig_ValidateEnabledRequiresURL(t *testing.T) {
	cfg := NtfyConfig{Enabled: true, Topic: "alerts", Priority: "default"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when server_url is empty")
	}
}

func TestNtfyConfig_ValidateEnabledRequiresTopic(t *testing.T) {
	cfg := NtfyConfig{Enabled: true, ServerURL: "https://ntfy.sh", Priority: "default"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when topic is empty")
	}
}

func TestNtfyConfig_ValidateInvalidPriority(t *testing.T) {
	cfg := NtfyConfig{
		Enabled:   true,
		ServerURL: "https://ntfy.sh",
		Topic:     "alerts",
		Priority:  "extreme",
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid priority")
	}
}

func TestNtfyConfig_ValidateValidConfig(t *testing.T) {
	cfg := NtfyConfig{
		Enabled:   true,
		ServerURL: "https://ntfy.sh",
		Topic:     "portwatch",
		Priority:  "high",
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
