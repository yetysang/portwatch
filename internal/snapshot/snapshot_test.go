package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
	"github.com/user/portwatch/internal/snapshot"
)

func tempStore(t *testing.T) *snapshot.Store {
	t.Helper()
	dir := t.TempDir()
	return snapshot.NewStore(filepath.Join(dir, "portwatch", "snapshot.json"))
}

func sampleBindings() []ports.Binding {
	return []ports.Binding{
		{Proto: "tcp", Addr: "0.0.0.0", Port: 8080, PID: 1234, Process: "nginx"},
		{Proto: "tcp", Addr: "127.0.0.1", Port: 5432, PID: 5678, Process: "postgres"},
	}
}

func TestStore_ExistsReturnsFalseWhenMissing(t *testing.T) {
	st := tempStore(t)
	if st.Exists() {
		t.Fatal("expected Exists() == false for new store")
	}
}

func TestStore_SaveAndLoad(t *testing.T) {
	st := tempStore(t)
	bindings := sampleBindings()

	if err := st.Save(bindings); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	snap, err := st.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if len(snap.Bindings) != len(bindings) {
		t.Fatalf("expected %d bindings, got %d", len(bindings), len(snap.Bindings))
	}
	if snap.Bindings[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", snap.Bindings[0].Port)
	}
	if snap.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}
}

func TestStore_ExistsAfterSave(t *testing.T) {
	st := tempStore(t)
	_ = st.Save(sampleBindings())
	if !st.Exists() {
		t.Fatal("expected Exists() == true after Save()")
	}
}

func TestStore_LoadEmptyWhenMissing(t *testing.T) {
	st := tempStore(t)
	snap, err := st.Load()
	if err != nil {
		t.Fatalf("Load() on missing file should not error: %v", err)
	}
	if len(snap.Bindings) != 0 {
		t.Errorf("expected 0 bindings, got %d", len(snap.Bindings))
	}
	if !snap.CapturedAt.Equal(time.Time{}) {
		t.Error("CapturedAt should be zero for empty snapshot")
	}
}

func TestStore_OverwritesPreviousSnapshot(t *testing.T) {
	st := tempStore(t)
	_ = st.Save(sampleBindings())

	newBindings := []ports.Binding{
		{Proto: "udp", Addr: "0.0.0.0", Port: 53, PID: 999, Process: "dnsmasq"},
	}
	_ = st.Save(newBindings)

	snap, err := st.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if len(snap.Bindings) != 1 {
		t.Fatalf("expected 1 binding after overwrite, got %d", len(snap.Bindings))
	}
	if snap.Bindings[0].Proto != "udp" {
		t.Errorf("expected proto udp, got %s", snap.Bindings[0].Proto)
	}
}

func TestStore_LoadCorruptFileReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	_ = os.WriteFile(path, []byte("not valid json{"), 0o644)

	st := snapshot.NewStore(path)
	_, err := st.Load()
	if err == nil {
		t.Fatal("expected error loading corrupt snapshot")
	}
}
