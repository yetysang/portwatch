package config

import "testing"

func TestDefaultDiscordConfig_Values(t *testing.T) {
	cfg := DefaultDiscordConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.WebhookURL != "" {
		t.Errorf("expected empty WebhookURL, got %q", cfg.WebhookURL)
	}
}

func TestDiscordConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := DiscordConfig{Enabled: false, WebhookURL: ""}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestDiscordConfig_ValidateEnabledRequiresURL(t *testing.T) {
	cfg := DiscordConfig{Enabled: true, WebhookURL: ""}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error when enabled without webhook_url")
	}
}

func TestDiscordConfig_ValidateEnabledWithURL(t *testing.T) {
	cfg := DiscordConfig{
		Enabled:    true,
		WebhookURL: "https://discord.com/api/webhooks/123/abc",
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestDiscordConfig_ValidateErrorMessage(t *testing.T) {
	cfg := DiscordConfig{Enabled: true, WebhookURL: ""}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error")
	}
	got := err.Error()
	want := "discord: webhook_url must be set when enabled"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
