// Package alert provides handlers for routing and dispatching port-change
// notifications to various backends.
//
// # Rollup Handler
//
// The RollupHandler batches multiple change events that arrive within a
// configurable time window into a single forwarded call. This reduces noise
// when many ports change state simultaneously — for example, during a service
// restart that closes and reopens dozens of sockets in quick succession.
//
// # Configuration
//
// Rollup behaviour is controlled by [config.RollupConfig]:
//
//	type RollupConfig struct {
//	    Enabled  bool
//	    Window   time.Duration // how long to accumulate changes
//	    MaxBatch int           // flush early if batch reaches this size
//	}
//
// When Enabled is false the handler forwards every call immediately without
// buffering, preserving the original latency characteristics.
//
// # Usage
//
//	handler := alert.NewRollupHandler(cfg.Rollup, next)
//
// The returned handler satisfies [alert.Handler] and can be composed with any
// other handler in the pipeline (e.g. wrapped by a [NewRetryHandler] or a
// [NewCircuitBreakerHandler]).
//
// # Flush semantics
//
// Changes are accumulated in an internal buffer. The buffer is flushed to the
// downstream handler when either:
//
//   - the rollup window elapses since the first buffered event, or
//   - the number of buffered changes reaches MaxBatch.
//
// A background goroutine drives the window-based flush. Call Drain on the
// handler during shutdown to flush any remaining buffered events before the
// process exits.
package alert
