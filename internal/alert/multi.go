package alert

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/monitor"
)

// MultiHandler fans out Change slices to multiple Handler implementations.
// All handlers are invoked regardless of individual errors; a combined error
// is returned if any handler fails.
type MultiHandler struct {
	handlers []Handler
}

// Handler is the interface implemented by all alert backends.
type Handler interface {
	Handle(changes []monitor.Change) error
	Drain() error
}

// NewMultiHandler returns a MultiHandler that delegates to each provided handler.
func NewMultiHandler(handlers ...Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

// Handle calls Handle on every registered handler and collects errors.
func (m *MultiHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	return m.collect(func(h Handler) error {
		return h.Handle(changes)
	})
}

// Drain calls Drain on every registered handler and collects errors.
func (m *MultiHandler) Drain() error {
	return m.collect(func(h Handler) error {
		return h.Drain()
	})
}

// Len returns the number of registered handlers.
func (m *MultiHandler) Len() int {
	return len(m.handlers)
}

func (m *MultiHandler) collect(fn func(Handler) error) error {
	var errs []string
	for _, h := range m.handlers {
		if err := fn(h); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("multi handler errors: %s", strings.Join(errs, "; "))
	}
	return nil
}
