package ports

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnricher_EnrichZeroInode(t *testing.T) {
	e := NewEnricher(nil)
	b := Binding{Port: 80, Proto: "tcp", Inode: 0}
	eb := e.Enrich(b)

	if eb.PID != 0 {
		t.Errorf("expected PID 0 for zero inode, got %d", eb.PID)
	}
	if eb.ProcessName != "" {
		t.Errorf("expected empty process name for zero inode, got %q", eb.ProcessName)
	}
	if eb.Port != 80 {
		t.Errorf("expected port 80, got %d", eb.Port)
	}
}

func TestEnricher_EnrichAll_PreservesOrder(t *testing.T) {
	e := NewEnricher(nil)
	bindings := []Binding{
		{Port: 22, Proto: "tcp"},
		{Port: 443, Proto: "tcp"},
		{Port: 8080, Proto: "tcp"},
	}

	result := e.EnrichAll(bindings)

	if len(result) != len(bindings) {
		t.Fatalf("expected %d enriched bindings, got %d", len(bindings), len(result))
	}
	for i, eb := range result {
		if eb.Port != bindings[i].Port {
			t.Errorf("index %d: expected port %d, got %d", i, bindings[i].Port, eb.Port)
		}
	}
}

func TestEnricher_FindProcess_FakeProcFS(t *testing.T) {
	root := t.TempDir()

	// create a fake /proc/1234/fd/3 -> socket:[9999]
	pidDir := filepath.Join(root, "1234")
	fdDir := filepath.Join(pidDir, "fd")
	if err := os.MkdirAll(fdDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// write comm file
	if err := os.WriteFile(filepath.Join(pidDir, "comm"), []byte("myapp\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// create symlink for socket inode
	symlink := filepath.Join(fdDir, "3")
	if err := os.Symlink("socket:[9999]", symlink); err != nil {
		t.Fatal(err)
	}

	e := &Enricher{procRoot: root}
	pid, name := e.findProcess(9999)

	if pid != 1234 {
		t.Errorf("expected PID 1234, got %d", pid)
	}
	if name != "myapp" {
		t.Errorf("expected process name 'myapp', got %q", name)
	}
}

func TestEnricher_FindProcess_MissingInode(t *testing.T) {
	root := t.TempDir()

	pidDir := filepath.Join(root, "42")
	fdDir := filepath.Join(pidDir, "fd")
	if err := os.MkdirAll(fdDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("socket:[1111]", filepath.Join(fdDir, "0")); err != nil {
		t.Fatal(err)
	}

	e := &Enricher{procRoot: root}
	pid, name := e.findProcess(9999) // different inode

	if pid != 0 {
		t.Errorf("expected PID 0, got %d", pid)
	}
	if name != "" {
		t.Errorf("expected empty name, got %q", name)
	}
}
