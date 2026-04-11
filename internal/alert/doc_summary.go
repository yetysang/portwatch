// Package alert provides handler implementations for dispatching port-change
// notifications through multiple output channels.
//
// # Summary Handler
//
// SummaryHandler accumulates Change events in memory and periodically emits a
// grouped, human-readable digest rather than one line per event. This is
// useful when port churn is high and you only need a periodic overview.
//
// Usage:
//
//	h := alert.NewSummaryHandler(os.Stdout, "[portwatch] ", 5*time.Minute)
//	// inside your tick loop:
//	_ = h.Handle(changes)
//	// on shutdown, flush any remaining buffered events:
//	_ = h.Flush()
//
// Configuration is managed through SummaryConfig, which can be embedded in
// the top-level application config and validated with SummaryConfig.Validate.
package alert
