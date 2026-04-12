// Package alert provides alerting handlers for portwatch.
//
// OpsGenie Handler
//
// NewOpsGenieHandler creates an alert handler that sends notifications to
// OpsGenie via its Alerts API. Each batch of changes results in one API
// call, with a summary message derived from the first change in the batch.
//
// Usage:
//
//	h := alert.NewOpsGenieHandler("https://api.opsgenie.com/v2/alerts", "<api-key>")
//
// The handler sets the Authorization header using the provided API key and
// posts a JSON payload conforming to the OpsGenie Create Alert API.
// A non-2xx response is treated as an error.
//
// Drain is a no-op for this handler.
package alert
