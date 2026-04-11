// Package alert provides pluggable alert handlers for portwatch.
//
// # Email Handler
//
// NewEmailHandler sends an SMTP email for each batch of port-binding changes.
// It requires a populated EmailConfig struct:
//
//	cfg := alert.EmailConfig{
//	    SMTPAddr: "smtp.example.com:587",
//	    From:     "portwatch@example.com",
//	    To:       []string{"ops@example.com"},
//	    Subject:  "Port binding alert",          // optional
//	    Auth:     smtp.PlainAuth("", user, pass, host), // optional
//	}
//	h := alert.NewEmailHandler(cfg)
//
// The handler skips sending when the changes slice is empty, making it
// safe to wrap with ThrottleHandler or DedupHandler.
//
// The email body is rendered from a plain-text template listing each
// change with its kind, protocol, port, optional hostname, and process
// information when available.
package alert
