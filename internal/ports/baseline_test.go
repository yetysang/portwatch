package ports

import (
	"os"
	"path/filepath"
	"testing"
)

func tempBaseline(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "baseline.json")
}

func TestNewBaseline_MissingFile(t *testing.T) {
	b, err := NewBaseline(tempBaseline(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Len() != 0 {
		t.Errorf("expected empty baseline, got %d entries", b.Len())
	}
}

func TestNewBaseline_InvalidJSON(t *testing.T) {
	path := tempBaseline(t)
	if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := NewBaseline(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestBaseline_AddAndContains(t *testing.T) {
	b, err := NewBaseline(tempBaseline(t))
	if err != nil {
		t.Fatal(err)
	}
	entry := BaselineEntry{Proto: "tcp", Addr: "0.0.0.0", Port: 8080, Process: "nginx"}
	if err := b.Add(entry); err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if !b.Contains("tcp", "0.0.0.0", 8080) {
		t.Error("expected baseline to contain tcp:0.0.0.0:8080")
	}
	if b.Contains("udp", "0.0.0.0", 8080) {
		t.Error("expected baseline NOT to contain udp:0.0.0.0:8080")
	}
}

func TestBaseline_Persistence(t *testing.T) {
	path := tempBaseline(t)
	b1, _ := NewBaseline(path)
	_ = b1.Add(BaselineEntry{Proto: "tcp", Addr: "127.0.0.1", Port: 5432})

	b2, err := NewBaseline(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if !b2.Contains("tcp", "127.0.0.1", 5432) {
		t.Error("expected reloaded baseline to contain entry")
	}
}

func TestBaseline_LenAfterMultipleAdds(t *testing.T) {
	b, _ := NewBaseline(tempBaseline(t))
	_ = b.Add(BaselineEntry{Proto: "tcp", Addr: "0.0.0.0", Port: 80})
	_ = b.Add(BaselineEntry{Proto: "tcp", Addr: "0.0.0.0", Port: 443})
	_ = b.Add(BaselineEntry{Proto: "udp", Addr: "0.0.0.0", Port: 53})
	if b.Len() != 3 {
		t.Errorf("expected 3 entries, got %d", b.Len())
	}
}

func TestBaseline_DuplicateAddDoesNotGrow(t *testing.T) {
	b, _ := NewBaseline(tempBaseline(t))
	entry := BaselineEntry{Proto: "tcp", Addr: "0.0.0.0", Port: 80}
	_ = b.Add(entry)
	_ = b.Add(entry)
	if b.Len() != 1 {
		t.Errorf("expected 1 entry after duplicate add, got %d", b.Len())
	}
}
