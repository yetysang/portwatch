package config

import "testing"

func TestDefaultWebhookConfig_Values(t *testing.T) {
	c := DefaultWebhookConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false")
	}
	if c.URL != "" {
		t.Errorf("expected empty URL, got %q", c.URL)
	}
	if c.Timeout != 5 {
		t.Errorf("expected Timeout=5, got %d", c.Timeout)
	}
}

func TestWebhookConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	c := WebhookConfig{Enabled: false, URL: "", Timeout: 0}
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestWebhookConfig_ValidateEnabledRequiresURL(t *testing.T) {
	c := WebhookConfig{Enabled: true, URL: "", Timeout: 5}
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing URL")
	}
}

func TestWebhookConfig_ValidateEnabledWithURL(t *testing.T) {
	c := WebhookConfig{Enabled: true, URL: "http://example.com/hook", Timeout: 5}
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestWebhookConfig_ValidateTimeoutTooLow(t *testing.T) {
	c := WebhookConfig{Enabled: true, URL: "http://example.com/hook", Timeout: 0}
	if err := c.Validate(); err == nil {
		t.Error("expected error for zero timeout")
	}
}

func TestWebhookConfig_ValidateTimeoutTooHigh(t *testing.T) {
	c := WebhookConfig{Enabled: true, URL: "http://example.com/hook", Timeout: 61}
	err := c.Validate()
	if err == nil {
		t.Error("expected error for timeout > 60")
	}
}
