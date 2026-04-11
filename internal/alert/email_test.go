package alert

import (
	"net/smtp"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func emailChange(kind monitor.ChangeKind, port uint16) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Port:     port,
			Protocol: "tcp",
			Addr:     "0.0.0.0",
		},
	}
}

func TestEmailHandler_EmptyChangesNoSend(t *testing.T) {
	sent := false
	h := &emailHandler{
		cfg: EmailConfig{From: "a@b.com", To: []string{"c@d.com"}, SMTPAddr: "localhost:25"},
		send: func(_ string, _ smtp.Auth, _ string, _ []string, _ []byte) error {
			sent = true
			return nil
		},
	}
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sent {
		t.Error("expected no send on empty changes")
	}
}

func TestEmailHandler_SendsOnChange(t *testing.T) {
	var capturedMsg []byte
	h := &emailHandler{
		cfg: EmailConfig{
			From:     "alert@portwatch.local",
			To:       []string{"ops@example.com"},
			SMTPAddr: "smtp.example.com:587",
			Subject:  "test alert",
		},
		send: func(_ string, _ smtp.Auth, _ string, _ []string, msg []byte) error {
			capturedMsg = msg
			return nil
		},
	}

	changes := []monitor.Change{emailChange(monitor.Added, 8080)}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body := string(capturedMsg)
	if !strings.Contains(body, "Subject: test alert") {
		t.Errorf("expected subject in message, got:\n%s", body)
	}
	if !strings.Contains(body, "8080") {
		t.Errorf("expected port 8080 in message body, got:\n%s", body)
	}
	if !strings.Contains(body, "added") {
		t.Errorf("expected change kind in message body, got:\n%s", body)
	}
}

func TestEmailHandler_DefaultSubject(t *testing.T) {
	var capturedMsg []byte
	h := &emailHandler{
		cfg: EmailConfig{From: "a@b.com", To: []string{"c@d.com"}, SMTPAddr: "localhost:25"},
		send: func(_ string, _ smtp.Auth, _ string, _ []string, msg []byte) error {
			capturedMsg = msg
			return nil
		},
	}
	if err := h.Handle([]monitor.Change{emailChange(monitor.Removed, 443)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(capturedMsg), "portwatch: port binding alert") {
		t.Errorf("expected default subject, got:\n%s", string(capturedMsg))
	}
}

func TestEmailHandler_DrainIsNoop(t *testing.T) {
	h := NewEmailHandler(EmailConfig{})
	if err := h.Drain(); err != nil {
		t.Fatalf("Drain should be a no-op, got: %v", err)
	}
}
