package alert

import "github.com/user/portwatch/internal/monitor"

// Handler is the core interface every alert sink must satisfy.
type Handler interface {
	// Handle processes a slice of port-binding changes. Implementations must
	// be safe to call with an empty slice and should return nil in that case.
	Handle(changes []monitor.Change) error
}

// HandlerFunc is a function adapter for Handler.
type HandlerFunc func([]monitor.Change) error

// Handle implements Handler.
func (f HandlerFunc) Handle(changes []monitor.Change) error {
	return f(changes)
}

// NopHandler discards all changes without error. Useful as a placeholder or in
// tests.
type NopHandler struct{}

// Handle implements Handler.
func (NopHandler) Handle(_ []monitor.Change) error { return nil }
