// VictorOps Handler
//
// NewVictorOpsHandler creates an alert handler that posts notifications to a
// VictorOps REST monitoring endpoint. The webhook URL should include the
// routing key, e.g.:
//
//	https://alert.victorops.com/integrations/generic/20131114/alert/<api-key>/<routing-key>
//
// Payload fields:
//
//	- message_type:    "CRITICAL" for added bindings, "RECOVERY" for removed.
//	- entity_id:       proto/port string used to correlate alert and recovery.
//	- state_message:   human-readable description of the change.
//	- monitoring_tool: always "portwatch".
//
// Only the first change in a batch is used to set the message_type and
// entity_id; the full batch count is not currently summarised.
//
// Drain is a no-op for this handler.
package alert
