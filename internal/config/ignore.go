package config

// IgnoreSet provides O(1) lookup for ports that should be suppressed.
type IgnoreSet struct {
	ports map[int]struct{}
}

// NewIgnoreSet builds an IgnoreSet from a slice of port numbers.
func NewIgnoreSet(ports []int) *IgnoreSet {
	s := &IgnoreSet{ports: make(map[int]struct{}, len(ports))}
	for _, p := range ports {
		s.ports[p] = struct{}{}
	}
	return s
}

// Contains reports whether port is in the ignore list.
func (s *IgnoreSet) Contains(port int) bool {
	_, ok := s.ports[port]
	return ok
}

// Len returns the number of ignored ports.
func (s *IgnoreSet) Len() int {
	return len(s.ports)
}
