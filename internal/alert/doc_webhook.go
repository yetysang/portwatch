// Package alert provides alerting handlers for portwatch.
//
// # Webhook Handler
//
// NewWebhookHandler sends a JSON POST request to a configured URL whenever
// port binding changes are detected.
//
// # Configuration
//
//	type WebhookConfig struct {
//		URL     string        // destination endpoint (required)
//		Timeout time.Duration // HTTP client timeout (default: 5s)
//		Secret  string        // optional HMAC-SHA256 signing secret
//	}
//
// # Payload
//
// Each request body is a JSON array of Change objects. If Secret is set,
// the handler adds an X-PortWatch-Signature header containing the hex-encoded
// HMAC-SHA256 of the raw body, prefixed with "sha256=".
//
// # Example
//
//	h := alert.NewWebhookHandler(alert.WebhookConfig{
//		URL:    "https://example.com/hooks/portwatch",
//		Secret: os.Getenv("HOOK_SECRET"),
//	})
package alert
