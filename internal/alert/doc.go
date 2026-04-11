// Package alert provides composable handlers for dispatching port-change
// notifications to various sinks.
//
// # Handlers
//
// Each handler implements the Handler interface:
//
//	type Handler interface {
//		Handle(changes []monitor.Change) error
//	}
//
// Available implementations:
//
//   - [NewHandler]           – structured log handler (zerolog)
//   - [NewStdoutHandler]     – human-readable stdout lines
//   - [NewFileHandler]       – append JSON-lines to a file
//   - [NewWebhookHandler]    – HTTP POST JSON payload
//   - [NewMultiHandler]      – fan-out to multiple handlers
//   - [NewRateLimitedHandler]– per-change-key cooldown wrapper
//   - [NewThrottleHandler]   – whole-batch quiet-window wrapper
//
// # Composition
//
// Handlers are designed to be composed. A typical production setup:
//
//	multi := alert.NewMultiHandler(
//		alert.NewRateLimitedHandler(alert.NewStdoutHandler(os.Stdout), cfg),
//		alert.NewThrottleHandler(alert.NewWebhookHandler(webhookURL), 30*time.Second),
//	)
package alert
