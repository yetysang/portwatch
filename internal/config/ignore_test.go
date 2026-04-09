package config

import "testing"

func TestNewIgnoreSet_Empty(t *testing.T) {
	s := NewIgnoreSet(nil)
	if s.Len() != 0 {
		t.Errorf("expected 0 entries, got %d", s.Len())
	}
	if s.Contains(80) {
		t.Error("empty set should not contain 80")
	}
}

func TestNewIgnoreSet_Contains(t *testing.T) {
	s := NewIgnoreSet([]int{22, 80, 443})

	for _, p := range []int{22, 80, 443} {
		if !s.Contains(p) {
			t.Errorf("expected set to contain port %d", p)
		}
	}
}

func TestNewIgnoreSet_NotContains(t *testing.T) {
	s := NewIgnoreSet([]int{22, 80})
	if s.Contains(8080) {
		t.Error("set should not contain 8080")
	}
}

func TestNewIgnoreSet_Len(t *testing.T) {
	s := NewIgnoreSet([]int{1, 2, 3, 4, 5})
	if s.Len() != 5 {
		t.Errorf("expected 5, got %d", s.Len())
	}
}

func TestNewIgnoreSet_Dedup(t *testing.T) {
	// map naturally deduplicates
	s := NewIgnoreSet([]int{80, 80, 80})
	if s.Len() != 1 {
		t.Errorf("expected 1 after dedup, got %d", s.Len())
	}
}
