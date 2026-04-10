package ports

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func makeBindings(ports ...int) []Binding {
	out := make([]Binding, 0, len(ports))
	for _, p := range ports {
		out = append(out, Binding{Proto: "tcp", Addr: "0.0.0.0", Port: p, PID: 1})
	}
	return out
}

func TestFilter_NilIgnorePassesAll(t *testing.T) {
	f := NewFilter(nil)
	bindings := makeBindings(80, 443, 8080)
	got := f.Apply(bindings)
	if len(got) != len(bindings) {
		t.Fatalf("expected %d bindings, got %d", len(bindings), len(got))
	}
}

func TestFilter_EmptyIgnorePassesAll(t *testing.T) {
	ignore := config.NewIgnoreSet(nil)
	f := NewFilter(ignore)
	bindings := makeBindings(80, 443)
	got := f.Apply(bindings)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilter_RemovesIgnoredPorts(t *testing.T) {
	ignore := config.NewIgnoreSet([]int{80, 443})
	f := NewFilter(ignore)
	bindings := makeBindings(80, 443, 8080, 9090)
	got := f.Apply(bindings)
	if len(got) != 2 {
		t.Fatalf("expected 2 bindings after filter, got %d", len(got))
	}
	for _, b := range got {
		if b.Port == 80 || b.Port == 443 {
			t.Errorf("port %d should have been filtered out", b.Port)
		}
	}
}

func TestFilter_ApplyToMap(t *testing.T) {
	ignore := config.NewIgnoreSet([]int{22})
	f := NewFilter(ignore)
	m := map[string]Binding{
		"tcp:22":   {Proto: "tcp", Port: 22},
		"tcp:8080": {Proto: "tcp", Port: 8080},
	}
	got := f.ApplyToMap(m)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if _, ok := got["tcp:22"]; ok {
		t.Error("port 22 should have been filtered from map")
	}
}

func TestFilter_ApplyToMap_NilIgnore(t *testing.T) {
	f := NewFilter(nil)
	m := map[string]Binding{
		"tcp:80": {Proto: "tcp", Port: 80},
	}
	got := f.ApplyToMap(m)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
}
