package config

import (
	"testing"
	"time"
)

func TestDefaultMattermostConfig_Values(t *testing.T) {
	c := DefaultMattermostConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false")
	}
	if c.Username != "portwatch" {
		t.Errorf("expected username portwatch, got %q", c.Username)
	}
	if c.Timeout != 10*time.Second {
		t.Errorf("expected timeout 10s, got %v", c.Timeout)
	}
}

func TestMattermostConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	c := DefaultMattermostConfig()
	c.Enabled = false
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestMattermostConfig_ValidateEnabledRequiresURL(t *testing.T) {
	c := DefaultMattermostConfig()
	c.Enabled = true
	c.Channel = "#alerts"
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing URL")
	}
}

func TestMattermostConfig_ValidateEnabledRequiresChannel(t *testing.T) {
	c := DefaultMattermostConfig()
	c.Enabled = true
	c.URL = "https://mattermost.example.com/hooks/abc"
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing channel")
	}
}

func TestMattermostConfig_ValidateEnabledWithBothFields(t *testing.T) {
	c := DefaultMattermostConfig()
	c.Enabled = true
	c.URL = "https://mattermost.example.com/hooks/abc"
	c.Channel = "#alerts"
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
