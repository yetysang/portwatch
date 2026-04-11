package alert

import (
	"sync"

	"github.com/iamcalledrob/portwatch/internal/monitor"
)

// DedupHandler wraps a Handler and suppresses duplicate changes within a
// single Handle call. Two changes are considered duplicates if they share
// the same key (protocol + address + port + kind).
type DedupHandler struct {
	mu   sync.Mutex
	seen map[string]struct{}
	next Handler
}

// NewDedupHandler returns a DedupHandler that forwards unique changes to next.
func NewDedupHandler(next Handler) *DedupHandler {
	return &DedupHandler{
		next: next,
		seen: make(map[string]struct{}),
	}
}

// Handle filters out duplicate changes within the slice before forwarding.
func (d *DedupHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	unique := changes[:0:0]
	local := make(map[string]struct{}, len(changes))

	for _, c := range changes {
		k := dedupKey(c)
		if _, seen := local[k]; seen {
			continue
		}
		local[k] = struct{}{}
		unique = append(unique, c)
	}

	if len(unique) == 0 {
		return nil
	}
	return d.next.Handle(unique)
}

// Drain flushes any buffered state in the underlying handler if it implements Drainer.
func (d *DedupHandler) Drain() error {
	if dr, ok := d.next.(interface{ Drain() error }); ok {
		return dr.Drain()
	}
	return nil
}

func dedupKey(c monitor.Change) string {
	return changeKey(c)
}
