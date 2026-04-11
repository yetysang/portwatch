package ports

import (
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func TestBindingsToMap_Empty(t *testing.T) {
	m := BindingsToMap(nil)
	if len(m) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(m))
	}
}

func TestBindingsToMap_UniqueKeys(t *testing.T) {
	bindings := []Binding{
		{Proto: "tcp", Addr: "0.0.0.0", Port: 80},
		{Proto: "tcp", Addr: "0.0.0.0", Port: 443},
		{Proto: "udp", Addr: "0.0.0.0", Port: 53},
	}
	m := BindingsToMap(bindings)
	if len(m) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(m))
	}
}

func TestDiff_NoChanges(t *testing.T) {
	bindings := []Binding{
		{Proto: "tcp", Addr: "0.0.0.0", Port: 8080},
	}
	prev := BindingsToMap(bindings)
	curr := BindingsToMap(bindings)
	changes := Diff(prev, curr)
	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %d", len(changes))
	}
}

func TestDiff_DetectsAdded(t *testing.T) {
	prev := BindingsToMap(nil)
	curr := BindingsToMap([]Binding{
		{Proto: "tcp", Addr: "127.0.0.1", Port: 9000, Process: "myapp"},
	})
	changes := Diff(prev, curr)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != monitor.ChangeAdded {
		t.Errorf("expected ChangeAdded, got %v", changes[0].Type)
	}
	if changes[0].Port != 9000 {
		t.Errorf("expected port 9000, got %d", changes[0].Port)
	}
	if changes[0].Process != "myapp" {
		t.Errorf("expected process myapp, got %q", changes[0].Process)
	}
}

func TestDiff_DetectsRemoved(t *testing.T) {
	prev := BindingsToMap([]Binding{
		{Proto: "tcp", Addr: "0.0.0.0", Port: 3000},
	})
	curr := BindingsToMap(nil)
	changes := Diff(prev, curr)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != monitor.ChangeRemoved {
		t.Errorf("expected ChangeRemoved, got %v", changes[0].Type)
	}
	if changes[0].Port != 3000 {
		t.Errorf("expected port 3000, got %d", changes[0].Port)
	}
}

func TestDiff_AddedAndRemoved(t *testing.T) {
	prev := BindingsToMap([]Binding{
		{Proto: "tcp", Addr: "0.0.0.0", Port: 80},
	})
	curr := BindingsToMap([]Binding{
		{Proto: "tcp", Addr: "0.0.0.0", Port: 8080},
	})
	changes := Diff(prev, curr)
	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}
	var added, removed int
	for _, c := range changes {
		switch c.Type {
		case monitor.ChangeAdded:
			added++
		case monitor.ChangeRemoved:
			removed++
		}
	}
	if added != 1 || removed != 1 {
		t.Errorf("expected 1 added and 1 removed, got %d/%d", added, removed)
	}
}

func TestDiff_BothEmpty(t *testing.T) {
	prev := BindingsToMap(nil)
	curr := BindingsToMap(nil)
	changes := Diff(prev, curr)
	if len(changes) != 0 {
		t.Fatalf("expected no changes for two empty maps, got %d", len(changes))
	}
}
