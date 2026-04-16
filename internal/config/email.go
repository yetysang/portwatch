package config

import "fmt"

// EmailConfig holds configuration for the SMTP email alert handler.
type EmailConfig struct {
	Enabled    bool     `toml:"enabled"`
	SMTPHost   string   `toml:"smtp_host"`
	SMTPPort   int      `toml:"smtp_port"`
	Username   string   `toml:"username"`
	Password   string   `toml:"password"`
	From       string   `toml:"from"`
	To         []string `toml:"to"`
	Subject    string   `toml:"subject"`
	UseTLS     bool     `toml:"use_tls"`
}

// DefaultEmailConfig returns an EmailConfig with sensible defaults.
func DefaultEmailConfig() EmailConfig {
	return EmailConfig{
		Enabled:  false,
		SMTPPort: 587,
		Subject:  "[portwatch] Port binding change detected",
		UseTLS:   true,
	}
}

// Validate returns an error if the EmailConfig is invalid.
func (c EmailConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.SMTPHost == "" {
		return fmt.Errorf("email: smtp_host is required when enabled")
	}
	if c.From == "" {
		return fmt.Errorf("email: from address is required when enabled")
	}
	if len(c.To) == 0 {
		return fmt.Errorf("email: at least one recipient in 'to' is required when enabled")
	}
	if c.SMTPPort < 1 || c.SMTPPort > 65535 {
		return fmt.Errorf("email: smtp_port must be between 1 and 65535")
	}
	return nil
}
