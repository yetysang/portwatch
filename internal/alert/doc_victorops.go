// Package alert provides alert handlers for portwatch.
//
// # VictorOps (Splunk On-Call) Handler
//
// NewVictorOpsHandler sends port-change alerts to VictorOps (now known as
// Splunk On-Call) using the REST endpoint integration.
//
// Each binding change is mapped to a VictorOps message_type:
//
//	- Added bindings  → CRITICAL
//	- Removed bindings → RECOVERY
//
// Configuration:
//
//	victorops:
//	  enabled: true
//	  url: "https://alert.victorops.com/integrations/generic/20131114/alert"
//	  routing_key: "your-routing-key"
//	  timeout: 10s
//
// The routing_key is appended to the URL as a path segment when sending
// the alert, following the VictorOps REST endpoint convention.
package alert
