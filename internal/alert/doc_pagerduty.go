// Package alert provides alerting handlers for portwatch.
//
// # OpsGenie Handler
//
// NewOpsGenieHandler sends alerts to OpsGenie when port bindings change.
// It uses the OpsGenie Alert API v2 to create and close alerts.
//
// Configuration:
//
//	type OpsGenieConfig struct {
//	    APIKey   string // OpsGenie API key (required)
//	    Team     string // OpsGenie team to route alerts to (optional)
//	    BaseURL  string // override for testing; defaults to https://api.opsgenie.com
//	}
//
// Each port-binding change produces one OpsGenie alert with:
//   - alias derived from protocol + port (for deduplication)
//   - priority P3 for added bindings, P5 for removed
//   - tags: ["portwatch", "<proto>", "<added|removed>"]
//
// Drain is a no-op; OpsGenie manages its own alert lifecycle.
package alert
