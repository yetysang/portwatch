// Package alert provides alert handlers for portwatch.
//
// # Zenduty Handler
//
// ZendutyHandler sends port binding change events to the Zenduty
// incident management platform via its Events API.
//
// Each change (added or removed binding) is delivered as a separate
// Zenduty alert event. The alert_type field controls incident severity
// and must be one of "info", "warning", or "critical".
//
// Configuration example (portwatch.toml):
//
//	[zenduty]
//	enabled        = true
//	api_key        = "your-zenduty-api-key"
//	service_id     = "your-service-uuid"
//	integration_id = "your-integration-uuid"
//	alert_type     = "warning"
//
// The handler is a no-op when Enabled is false.
package alert
