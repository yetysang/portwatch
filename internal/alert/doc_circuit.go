// Package alert provides handlers for dispatching port-change notifications
// to various backends.
//
// # Circuit Breaker Handler
//
// NewCircuitBreakerHandler wraps another Handler and prevents cascading
// failures by temporarily disabling forwarding when the downstream handler
// returns too many consecutive errors.
//
// The circuit breaker operates in three states:
//
//   - Closed (normal): all calls are forwarded to the inner handler.
//   - Open (tripped):  calls are dropped immediately; no requests are made
//     to the downstream backend. The circuit stays open for a
//     configurable cool-down window.
//   - Half-open (probing): after the cool-down expires one call is allowed
//     through. A success resets the circuit to Closed; a failure returns
//     it to Open for another full cool-down period.
//
// # Configuration
//
//	cb := alert.NewCircuitBreakerHandler(inner, alert.CircuitBreakerConfig{
//	    Threshold:  5,               // consecutive failures before tripping
//	    CoolDown:   30 * time.Second, // how long to stay open
//	})
//
// The zero value of CircuitBreakerConfig is not valid; always supply at
// least a positive Threshold and a non-zero CoolDown.
//
// # Observability
//
// State transitions are logged at WARN level so that operators can detect
// flapping backends without being flooded with individual error messages.
package alert
