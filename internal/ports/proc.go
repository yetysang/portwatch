// Package ports provides utilities for scanning and filtering active port bindings.
package ports

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ProcInfo holds basic process metadata associated with a port binding.
type ProcInfo struct {
	PID  int
	Name string
}

// LookupProc attempts to resolve the process name for the given PID by reading
// /proc/<pid>/comm. Returns an empty ProcInfo if the lookup fails.
func LookupProc(pid int) ProcInfo {
	if pid <= 0 {
		return ProcInfo{PID: pid}
	}
	commPath := fmt.Sprintf("/proc/%d/comm", pid)
	data, err := os.ReadFile(commPath)
	if err != nil {
		return ProcInfo{PID: pid}
	}
	name := strings.TrimSpace(string(data))
	return ProcInfo{PID: pid, Name: name}
}

// ParseInode extracts the inode number from a /proc/net/tcp* line field.
// The inode is the 10th whitespace-separated field (0-indexed: 9).
func ParseInode(fields []string) (uint64, error) {
	if len(fields) < 10 {
		return 0, fmt.Errorf("too few fields: %d", len(fields))
	}
	inode, err := strconv.ParseUint(fields[9], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid inode %q: %w", fields[9], err)
	}
	return inode, nil
}

// String returns a human-readable representation of ProcInfo.
// If the process name is known, it formats as "<name>(<pid>)".
// Otherwise it falls back to just the PID.
func (p ProcInfo) String() string {
	if p.Name != "" {
		return fmt.Sprintf("%s(%d)", p.Name, p.PID)
	}
	return fmt.Sprintf("%d", p.PID)
}
