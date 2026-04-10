package ports

import "github.com/user/portwatch/internal/config"

// Filter wraps an IgnoreSet and provides port-level filtering
// for scanner results.
type Filter struct {
	ignore *config.IgnoreSet
}

// NewFilter creates a Filter backed by the given IgnoreSet.
func NewFilter(ignore *config.IgnoreSet) *Filter {
	return &Filter{ignore: ignore}
}

// Apply returns only those Bindings whose port is not in the ignore set.
func (f *Filter) Apply(bindings []Binding) []Binding {
	if f.ignore == nil || f.ignore.Len() == 0 {
		return bindings
	}
	out := make([]Binding, 0, len(bindings))
	for _, b := range bindings {
		if !f.ignore.Contains(b.Port) {
			out = append(out, b)
		}
	}
	return out
}

// ApplyToMap filters a map[string]Binding (keyed by any string) by port.
func (f *Filter) ApplyToMap(m map[string]Binding) map[string]Binding {
	if f.ignore == nil || f.ignore.Len() == 0 {
		return m
	}
	out := make(map[string]Binding, len(m))
	for k, b := range m {
		if !f.ignore.Contains(b.Port) {
			out[k] = b
		}
	}
	return out
}
