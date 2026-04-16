package alert

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/user/portwatch/internal/config"
)

// smtpSender is a function type for sending mail, injectable for testing.
type smtpSender func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

// buildEmailBody constructs a plain-text email body from a list of changes.
func buildEmailBody(cfg config.EmailConfig, lines []string) []byte {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("From: %s\r\n", cfg.From))
	sb.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(cfg.To, ", ")))
	sb.WriteString(fmt.Sprintf("Subject: %s\r\n", cfg.Subject))
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	sb.WriteString("\r\n")
	for _, l := range lines {
		sb.WriteString(l)
		sb.WriteString("\r\n")
	}
	return []byte(sb.String())
}

// dialTLS opens a TLS connection to the SMTP server.
func dialTLS(host string, port int) (*smtp.Client, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	tlsCfg := &tls.Config{ServerName: host}
	conn, err := tls.Dial("tcp", addr, tlsCfg)
	if err != nil {
		return nil, err
	}
	return smtp.NewClient(conn, host)
}

// dialPlain opens a plain TCP connection to the SMTP server.
func dialPlain(host string, port int) (*smtp.Client, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	return smtp.Dial(addr)
}

// resolveAuth returns smtp.Auth if credentials are provided.
func resolveAuth(cfg config.EmailConfig, host string) smtp.Auth {
	if cfg.Username == "" {
		return nil
	}
	return smtp.PlainAuth("", cfg.Username, cfg.Password, net.JoinHostPort(host, fmt.Sprintf("%d", cfg.SMTPPort)))
}
