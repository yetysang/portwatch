package alert

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func fileChange(action monitor.Action, port int) monitor.Change {
	return monitor.Change{
		Action: action,
		Binding: ports.Binding{
			Proto:   "tcp",
			Addr:    "0.0.0.0",
			Port:    port,
			PID:     1234,
			Process: "nginx",
		},
	}
}

func TestFileHandler_EmptyChangesNoWrite(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "alert-*.log")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	h, err := NewFileHandler(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()

	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, _ := os.Stat(tmp.Name())
	if info.Size() != 0 {
		t.Errorf("expected empty file, got %d bytes", info.Size())
	}
}

func TestFileHandler_WritesJSONLines(t *testing.T) {
	dir := t.TempDir()
	h, err := NewFileHandler(dir + "/out.log")
	if err != nil {
		t.Fatal(err)
	}

	changes := []monitor.Change{
		fileChange(monitor.Added, 80),
		fileChange(monitor.Removed, 443),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("Handle error: %v", err)
	}
	h.Close()

	f, _ := os.Open(dir + "/out.log")
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var records []fileRecord
	for scanner.Scan() {
		var rec fileRecord
		if err := json.Unmarshal(scanner.Bytes(), &rec); err != nil {
			t.Fatalf("invalid JSON line: %v", err)
		}
		records = append(records, rec)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].Port != 80 || records[0].Action != string(monitor.Added) {
		t.Errorf("unexpected first record: %+v", records[0])
	}
	if records[1].Port != 443 || records[1].Action != string(monitor.Removed) {
		t.Errorf("unexpected second record: %+v", records[1])
	}
}

func TestFileHandler_InvalidPath(t *testing.T) {
	_, err := NewFileHandler("/no/such/dir/out.log")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestFileHandler_AppendsBetweenCalls(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/append.log"

	for i := 0; i < 3; i++ {
		h, err := NewFileHandler(path)
		if err != nil {
			t.Fatal(err)
		}
		_ = h.Handle([]monitor.Change{fileChange(monitor.Added, 8080+i)})
		h.Close()
	}

	f, _ := os.Open(path)
	defer f.Close()
	count := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 appended lines, got %d", count)
	}
}
