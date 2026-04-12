// Package alert provides alerting handlers for portwatch.
//
// # Datadog Handler
//
// The Datadog handler posts port binding changes as events to the
// Datadog Events API (v1). Each change is submitted as a separate
// event with an appropriate alert type:
//
//   - "added" bindings → alert_type "warning"
//   - "removed" bindings → alert_type "info"
//
// Configuration fields (via config.DatadogConfig):
//
//	Enabled  bool   — enable or disable the handler
//	APIKey   string — Datadog API key (required when enabled)
//	Site     string — Datadog site, e.g. "datadoghq.com" (default)
//	Tags     []string — optional tags attached to every event
//
// Events are submitted to:
//
//	https://api.<site>/api/v1/events
//
// The handler skips submission when the changes slice is empty.
package alert
