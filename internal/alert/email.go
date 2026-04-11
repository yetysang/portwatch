package alert

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// EmailConfig holds SMTP configuration for the email alert handler.
type EmailConfig struct {
	SMTPAddr string // host:port
	From     string
	To       []string
	Subject  string
	Auth     smtp.Auth // nil for unauthenticated
}

var emailTmpl = template.Must(template.New("email").Parse(
	"Port binding changes detected at {{.Time}}:\n\n" +
		"{{range .Changes}}  [{{.Kind}}] {{.Protocol}} :{{.Port}}" +
		"{{if .Hostname}} ({{.Hostname}}){{end}}" +
		"{{if .Process}} — {{.Process}} (PID {{.PID}}){{end}}\n{{end}}",
))

type emailHandler struct {
	cfg  EmailConfig
	send func(addr string, auth smtp.Auth, from string, to []string, msg []byte) error
}

// NewEmailHandler returns a Handler that sends an SMTP email on each batch of changes.
func NewEmailHandler(cfg EmailConfig) Handler {
	return &emailHandler{cfg: cfg, send: smtp.SendMail}
}

func (h *emailHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	type data struct {
		Time    string
		Changes []monitor.Change
	}

	var body bytes.Buffer
	if err := emailTmpl.Execute(&body, data{
		Time:    time.Now().UTC().Format(time.RFC3339),
		Changes: changes,
	}); err != nil {
		return fmt.Errorf("email: render template: %w", err)
	}

	subject := h.cfg.Subject
	if subject == "" {
		subject = "portwatch: port binding alert"
	}

	var msg bytes.Buffer
	fmt.Fprintf(&msg, "From: %s\r\n", h.cfg.From)
	fmt.Fprintf(&msg, "To: %s\r\n", strings.Join(h.cfg.To, ", "))
	fmt.Fprintf(&msg, "Subject: %s\r\n", subject)
	fmt.Fprintf(&msg, "Content-Type: text/plain; charset=utf-8\r\n\r\n")
	msg.Write(body.Bytes())

	if err := h.send(h.cfg.SMTPAddr, h.cfg.Auth, h.cfg.From, h.cfg.To, msg.Bytes()); err != nil {
		return fmt.Errorf("email: send: %w", err)
	}
	return nil
}

func (h *emailHandler) Drain() error { return nil }
