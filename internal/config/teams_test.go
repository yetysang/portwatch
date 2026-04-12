package config

import "testing"

func TestDefaultTeamsConfig_Values(t *testing.T) {
	cfg := DefaultTeamsConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.WebhookURL != "" {
		t.Errorf("expected empty WebhookURL, got %q", cfg.WebhookURL)
	}
}

func TestTeamsConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := TeamsConfig{Enabled: false, WebhookURL: ""}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestTeamsConfig_ValidateEnabledRequiresURL(t *testing.T) {
	cfg := TeamsConfig{Enabled: true, WebhookURL: ""}
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error when enabled with empty webhook_url")
	}
}

func TestTeamsConfig_ValidateEnabledWithURL(t *testing.T) {
	cfg := TeamsConfig{
		Enabled:    true,
		WebhookURL: "https://outlook.office.com/webhook/abc123",
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTeamsConfig_ValidateErrorMessage(t *testing.T) {
	cfg := TeamsConfig{Enabled: true, WebhookURL: ""}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error")
	}
	got := err.Error()
	want := "teams: webhook_url must not be empty when enabled"
	if got != want {
		t.Errorf("expected error %q, got %q", want, got)
	}
}
