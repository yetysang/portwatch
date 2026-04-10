package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// Snapshot holds a persisted set of port bindings at a point in time.
type Snapshot struct {
	CapturedAt time.Time        `json:"captured_at"`
	Bindings   []ports.Binding  `json:"bindings"`
}

// Store handles reading and writing snapshots to disk.
type Store struct {
	path string
}

// NewStore creates a Store that persists snapshots at the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the current bindings to disk as a JSON snapshot.
func (s *Store) Save(bindings []ports.Binding) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	snap := Snapshot{
		CapturedAt: time.Now().UTC(),
		Bindings:   bindings,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Load reads the most recent snapshot from disk.
// Returns an empty Snapshot if the file does not exist.
func (s *Store) Load() (Snapshot, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return Snapshot{}, nil
		}
		return Snapshot{}, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}

// Exists reports whether a snapshot file is present on disk.
func (s *Store) Exists() bool {
	_, err := os.Stat(s.path)
	return err == nil
}
