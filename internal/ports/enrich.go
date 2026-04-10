// Package ports provides utilities for scanning and enriching port binding data.
package ports

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// EnrichedBinding extends a Binding with process and service metadata.
type EnrichedBinding struct {
	Binding
	PID         int
	ProcessName string
	ServiceName string
}

// Enricher attaches process and service information to raw Bindings.
type Enricher struct {
	resolver *Resolver
	procRoot string // override for testing; defaults to "/proc"
}

// NewEnricher returns an Enricher using the given Resolver.
func NewEnricher(r *Resolver) *Enricher {
	return &Enricher{resolver: r, procRoot: "/proc"}
}

// Enrich resolves process name and service name for a single Binding.
func (e *Enricher) Enrich(b Binding) EnrichedBinding {
	eb := EnrichedBinding{Binding: b}

	pid, name := e.findProcess(b.Inode)
	eb.PID = pid
	eb.ProcessName = name

	if e.resolver != nil {
		eb.ServiceName = e.resolver.LookupServiceName(b.Port, b.Proto)
	}

	return eb
}

// EnrichAll enriches a slice of Bindings.
func (e *Enricher) EnrichAll(bindings []Binding) []EnrichedBinding {
	out := make([]EnrichedBinding, 0, len(bindings))
	for _, b := range bindings {
		out = append(out, e.Enrich(b))
	}
	return out
}

// findProcess walks /proc to find the PID and process name owning inode.
func (e *Enricher) findProcess(inode uint64) (int, string) {
	if inode == 0 {
		return 0, ""
	}

	entries, err := os.ReadDir(e.procRoot)
	if err != nil {
		return 0, ""
	}

	target := fmt.Sprintf("socket:[%d]", inode)

	for _, entry := range entries {
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		fdDir := filepath.Join(e.procRoot, entry.Name(), "fd")
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}
		for _, fd := range fds {
			link, err := os.Readlink(filepath.Join(fdDir, fd.Name()))
			if err != nil {
				continue
			}
			if link == target {
				name := e.readProcessName(entry.Name())
				return pid, name
			}
		}
	}
	return 0, ""
}

// readProcessName reads the comm file for a given pid string.
func (e *Enricher) readProcessName(pidStr string) string {
	data, err := os.ReadFile(filepath.Join(e.procRoot, pidStr, "comm"))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
