// Package monitor detects changes in port bindings between scans.
package monitor

import (
	"fmt"

	"github.com/example/portwatch/internal/ports"
)

// Change describes a single port-binding event detected during a monitor tick.
type Change struct {
	// Binding is the port binding that was added or removed.
	Binding ports.Binding
	// Added is true when the binding appeared, false when it disappeared.
	Added bool
}

// String returns a human-readable summary of the change.
func (c Change) String() string {
	verb := "removed"
	if c.Added {
		verb = "added"
	}
	return fmt.Sprintf("[%s] %s:%d (%s)",
		verb, c.Binding.IP, c.Binding.Port, c.Binding.Protocol)
}

// FilterChanges returns only the changes for which keep returns true.
func FilterChanges(changes []Change, keep func(Change) bool) []Change {
	out := make([]Change, 0, len(changes))
	for _, c := range changes {
		if keep(c) {
			out = append(out, c)
		}
	}
	return out
}

// PartitionChanges splits changes into two slices: added and removed.
func PartitionChanges(changes []Change) (added, removed []Change) {
	for _, c := range changes {
		if c.Added {
			added = append(added, c)
		} else {
			removed = append(removed, c)
		}
	}
	return
}
