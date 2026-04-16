package alert

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func baseCfg() config.EmailConfig {
	cfg := config.DefaultEmailConfig()
	cfg.Enabled = true
	cfg.SMTPHost = "smtp.example.com"
	cfg.From = "alerts@example.com"
	cfg.To = []string{"ops@example.com", "dev@example.com"}
	cfg.Subject = "[portwatch] alert"
	return cfg
}

func TestBuildEmailBody_ContainsFrom(t *testing.T) {
	cfg := baseCfg()
	body := buildEmailBody(cfg, []string{"port 8080 added"})
	if !strings.Contains(string(body), "From: alerts@example.com") {
		t.Error("expected From header in body")
	}
}

func TestBuildEmailBody_ContainsTo(t *testing.T) {
	cfg := baseCfg()
	body := buildEmailBody(cfg, []string{"port 8080 added"})
	if !strings.Contains(string(body), "ops@example.com") {
		t.Error("expected To header in body")
	}
}

func TestBuildEmailBody_ContainsSubject(t *testing.T) {
	cfg := baseCfg()
	body := buildEmailBody(cfg, []string{"port 8080 added"})
	if !strings.Contains(string(body), "[portwatch] alert") {
		t.Error("expected Subject in body")
	}
}

func TestBuildEmailBody_ContainsLines(t *testing.T) {
	cfg := baseCfg()
	body := buildEmailBody(cfg, []string{"port 9090 removed", "port 8080 added"})
	s := string(body)
	if !strings.Contains(s, "port 9090 removed") {
		t.Error("expected first line in body")
	}
	if !strings.Contains(s, "port 8080 added") {
		t.Error("expected second line in body")
	}
}

func TestResolveAuth_NoCredentials(t *testing.T) {
	cfg := baseCfg()
	auth := resolveAuth(cfg, "smtp.example.com")
	if auth != nil {
		t.Error("expected nil auth when no username set")
	}
}

func TestResolveAuth_WithCredentials(t *testing.T) {
	cfg := baseCfg()
	cfg.Username = "user"
	cfg.Password = "pass"
	auth := resolveAuth(cfg, "smtp.example.com")
	if auth == nil {
		t.Error("expected non-nil auth when credentials set")
	}
}
