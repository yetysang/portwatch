// Package alert provides alerting handlers for portwatch.
//
// # Loki Handler
//
// The Loki handler pushes structured log streams to a Grafana Loki instance
// using the HTTP push API (/loki/api/v1/push).
//
// Each port-binding change is emitted as a separate Loki stream entry with
// labels derived from the change: job, proto, and kind (added/removed).
//
// # Configuration
//
//	[loki]
//	enabled   = true
//	url       = "http://localhost:3100"
//	job_label = "portwatch"
//
// The url field must point to the root of the Loki instance (no trailing
// path). The job_label is attached to every stream as the "job" label,
// making it easy to filter portwatch events in Grafana.
//
// # Usage
//
//	h := alert.NewLokiHandler(cfg.Loki)
//	// pass h to a MultiHandler or use directly in the tick loop
package alert
