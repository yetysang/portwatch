// Package ports provides utilities for scanning and comparing port bindings.
package ports

import (
	"fmt"

	"github.com/user/portwatch/internal/monitor"
)

// Diff computes the changes between two snapshots of port bindings.
// prev and curr are maps keyed by a stable binding identifier.
// Returns a slice of Change values describing added and removed bindings.
func Diff(prev, curr map[string]Binding) []monitor.Change {
	var changes []monitor.Change

	// Detect added bindings
	for key, b := range curr {
		if _, exists := prev[key]; !exists {
			changes = append(changes, monitor.Change{
				Type:    monitor.ChangeAdded,
				Proto:   b.Proto,
				Addr:    b.Addr,
				Port:    b.Port,
				PID:     b.PID,
				Process: b.Process,
				Service: b.Service,
			})
		}
	}

	// Detect removed bindings
	for key, b := range prev {
		if _, exists := curr[key]; !exists {
			changes = append(changes, monitor.Change{
				Type:    monitor.ChangeRemoved,
				Proto:   b.Proto,
				Addr:    b.Addr,
				Port:    b.Port,
				PID:     b.PID,
				Process: b.Process,
				Service: b.Service,
			})
		}
	}

	return changes
}

// BindingsToMap converts a slice of Binding into a map keyed by "proto:addr:port".
func BindingsToMap(bindings []Binding) map[string]Binding {
	m := make(map[string]Binding, len(bindings))
	for _, b := range bindings {
		key := bindingKey(b)
		m[key] = b
	}
	return m
}

// bindingKey returns a stable string identifier for a Binding in the form "proto:addr:port".
func bindingKey(b Binding) string {
	return fmt.Sprintf("%s:%s:%d", b.Proto, b.Addr, b.Port)
}
