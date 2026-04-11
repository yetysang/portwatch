package alert

import (
	"testing"

	"github.com/iamcalledrob/portwatch/internal/monitor"
	"github.com/iamcalledrob/portwatch/internal/ports"
)

func dedupChange(proto, addr string, port int, kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Proto:   proto,
			Addr:    addr,
			Port:    port,
		},
	}
}

func TestDedupHandler_EmptyChangesSkipped(t *testing.T) {
	collector := &collectingHandler{}
	h := NewDedupHandler(collector)

	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(collector.received) != 0 {
		t.Errorf("expected no calls, got %d", len(collector.received))
	}
}

func TestDedupHandler_UniqueChangesForwarded(t *testing.T) {
	collector := &collectingHandler{}
	h := NewDedupHandler(collector)

	changes := []monitor.Change{
		dedupChange("tcp", "0.0.0.0", 80, monitor.ChangeAdded),
		dedupChange("tcp", "0.0.0.0", 443, monitor.ChangeAdded),
	}

	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(collector.received) != 2 {
		t.Errorf("expected 2 changes, got %d", len(collector.received))
	}
}

func TestDedupHandler_DuplicatesWithinCallSuppressed(t *testing.T) {
	collector := &collectingHandler{}
	h := NewDedupHandler(collector)

	dup := dedupChange("tcp", "0.0.0.0", 8080, monitor.ChangeAdded)
	changes := []monitor.Change{dup, dup, dup}

	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(collector.received) != 1 {
		t.Errorf("expected 1 unique change, got %d", len(collector.received))
	}
}

func TestDedupHandler_DifferentKindsSamePortNotDeduped(t *testing.T) {
	collector := &collectingHandler{}
	h := NewDedupHandler(collector)

	changes := []monitor.Change{
		dedupChange("tcp", "127.0.0.1", 9090, monitor.ChangeAdded),
		dedupChange("tcp", "127.0.0.1", 9090, monitor.ChangeRemoved),
	}

	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(collector.received) != 2 {
		t.Errorf("expected 2 changes (different kinds), got %d", len(collector.received))
	}
}

// collectingHandler records all changes passed to Handle.
type collectingHandler struct {
	received []monitor.Change
}

func (c *collectingHandler) Handle(changes []monitor.Change) error {
	c.received = append(c.received, changes...)
	return nil
}
