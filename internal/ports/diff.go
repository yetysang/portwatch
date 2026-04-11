// Package ports provides utilities for scanning and comparing port bindings.
package ports

import "github.com/user/portwatch/internal/monitor"

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
		key := b.Proto + ":" + b.Addr + ":" + portStr(b.Port)
		m[key] = b
	}
	return m
}

func portStr(p int) string {
	if p == 0 {
		return "0"
	}
	return itoa(p)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}
