package config

import "testing"

func TestDefaultEmailConfig_Values(t *testing.T) {
	cfg := DefaultEmailConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false")
	}
	if cfg.SMTPPort != 587 {
		t.Errorf("expected SMTPPort 587, got %d", cfg.SMTPPort)
	}
	if cfg.Subject == "" {
		t.Error("expected non-empty default Subject")
	}
	if !cfg.UseTLS {
		t.Error("expected UseTLS to be true by default")
	}
}

func TestEmailConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	cfg := DefaultEmailConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestEmailConfig_ValidateEnabledRequiresSMTPHost(t *testing.T) {
	cfg := DefaultEmailConfig()
	cfg.Enabled = true
	cfg.From = "a@b.com"
	cfg.To = []string{"c@d.com"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when smtp_host is missing")
	}
}

func TestEmailConfig_ValidateEnabledRequiresFrom(t *testing.T) {
	cfg := DefaultEmailConfig()
	cfg.Enabled = true
	cfg.SMTPHost = "smtp.example.com"
	cfg.To = []string{"c@d.com"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when from is missing")
	}
}

func TestEmailConfig_ValidateEnabledRequiresTo(t *testing.T) {
	cfg := DefaultEmailConfig()
	cfg.Enabled = true
	cfg.SMTPHost = "smtp.example.com"
	cfg.From = "a@b.com"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when to list is empty")
	}
}

func TestEmailConfig_ValidateEnabledWithAllFields(t *testing.T) {
	cfg := DefaultEmailConfig()
	cfg.Enabled = true
	cfg.SMTPHost = "smtp.example.com"
	cfg.From = "a@b.com"
	cfg.To = []string{"c@d.com"}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestEmailConfig_ValidateInvalidPort(t *testing.T) {
	cfg := DefaultEmailConfig()
	cfg.Enabled = true
	cfg.SMTPHost = "smtp.example.com"
	cfg.From = "a@b.com"
	cfg.To = []string{"c@d.com"}
	cfg.SMTPPort = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid smtp_port")
	}
}
