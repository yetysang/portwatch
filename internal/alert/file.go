package alert

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// FileHandler writes change events as newline-delimited JSON to a file.
type FileHandler struct {
	path string
	f    *os.File
}

// NewFileHandler opens (or creates) the file at path for appending and returns
// a FileHandler. The caller is responsible for calling Close when done.
func NewFileHandler(path string) (*FileHandler, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("alert/file: open %q: %w", path, err)
	}
	return &FileHandler{path: path, f: f}, nil
}

// fileRecord is the JSON structure written for each change.
type fileRecord struct {
	Timestamp string `json:"timestamp"`
	Action    string `json:"action"`
	Proto     string `json:"proto"`
	Addr      string `json:"addr"`
	Port      int    `json:"port"`
	PID       int    `json:"pid,omitempty"`
	Process   string `json:"process,omitempty"`
}

// Handle writes each change in the slice as a JSON line to the log file.
func (h *FileHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	ts := time.Now().UTC().Format(time.RFC3339)
	enc := json.NewEncoder(h.f)
	for _, c := range changes {
		rec := fileRecord{
			Timestamp: ts,
			Action:    string(c.Action),
			Proto:     c.Binding.Proto,
			Addr:      c.Binding.Addr,
			Port:      c.Binding.Port,
			PID:       c.Binding.PID,
			Process:   c.Binding.Process,
		}
		if err := enc.Encode(rec); err != nil {
			return fmt.Errorf("alert/file: encode: %w", err)
		}
	}
	return nil
}

// Close closes the underlying file.
func (h *FileHandler) Close() error {
	return h.f.Close()
}
