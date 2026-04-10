package ports

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestLookupProc_InvalidPID(t *testing.T) {
	info := LookupProc(-1)
	if info.PID != -1 {
		t.Errorf("expected PID -1, got %d", info.PID)
	}
	if info.Name != "" {
		t.Errorf("expected empty name for invalid PID, got %q", info.Name)
	}
}

func TestLookupProc_ZeroPID(t *testing.T) {
	info := LookupProc(0)
	if info.Name != "" {
		t.Errorf("expected empty name for zero PID, got %q", info.Name)
	}
}

func TestLookupProc_CurrentProcess(t *testing.T) {
	pid := os.Getpid()
	info := LookupProc(pid)
	if info.PID != pid {
		t.Errorf("expected PID %d, got %d", pid, info.PID)
	}
	// Name may or may not resolve depending on environment; just ensure no panic.
}

func TestLookupProc_NonExistentPID(t *testing.T) {
	// Use an extremely large PID unlikely to exist.
	info := LookupProc(9999999)
	if info.Name != "" {
		t.Errorf("expected empty name for non-existent PID, got %q", info.Name)
	}
}

func TestParseInode_Valid(t *testing.T) {
	fields := make([]string, 10)
	fields[9] = "123456"
	inode, err := ParseInode(fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inode != 123456 {
		t.Errorf("expected inode 123456, got %d", inode)
	}
}

func TestParseInode_TooFewFields(t *testing.T) {
	_, err := ParseInode([]string{"a", "b"})
	if err == nil {
		t.Error("expected error for too few fields")
	}
}

func TestParseInode_InvalidValue(t *testing.T) {
	fields := make([]string, 10)
	fields[9] = "notanumber"
	_, err := ParseInode(fields)
	if err == nil {
		t.Error("expected error for non-numeric inode")
	}
}

func TestParseInode_LargeValue(t *testing.T) {
	fields := make([]string, 10)
	expected := uint64(4294967295)
	fields[9] = strconv.FormatUint(expected, 10)
	inode, err := ParseInode(fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inode != expected {
		t.Errorf("expected %d, got %d", expected, inode)
	}
	_ = fmt.Sprintf("inode=%d", inode) // suppress unused warning
}

func TestParseInode_ZeroInode(t *testing.T) {
	fields := make([]string, 10)
	fields[9] = "0"
	inode, err := ParseInode(fields)
	if err != nil {
		t.Fatalf("unexpected error for zero inode: %v", err)
	}
	if inode != 0 {
		t.Errorf("expected inode 0, got %d", inode)
	}
}
