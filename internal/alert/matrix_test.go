package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"portwatch/internal/monitor"
	"portwatch/internal/ports"
)

func matrixChange(action monitor.ChangeKind, port uint16, host string) monitor.Change {
	return monitor.Change{
		Kind: action,
		Binding: ports.Binding{
			Proto:    "tcp",
			Port:     port,
			Hostname: host,
			Addr:     "127.0.0.1",
		},
	}
}

func TestMatrixHandler_EmptyChangesNoRequest(t *testing.T) {
	requests := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewMatrixHandler(ts.URL, "token", "!room:example.org")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if requests != 0 {
		t.Errorf("expected 0 requests, got %d", requests)
	}
}

func TestMatrixHandler_PostsOnChange(t *testing.T) {
	var body []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"event_id":"$abc"}`))
	}))
	defer ts.Close()

	h := NewMatrixHandler(ts.URL, "mytoken", "!room:example.org")
	changes := []monitor.Change{
		matrixChange(monitor.ChangeAdded, 8080, "webserver"),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}
	msgtype, _ := payload["msgtype"].(string)
	if msgtype != "m.text" {
		t.Errorf("expected msgtype m.text, got %q", msgtype)
	}
	body2, _ := payload["body"].(string)
	if !strings.Contains(body2, "added") {
		t.Errorf("expected body to contain 'added', got %q", body2)
	}
}

func TestMatrixHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	h := NewMatrixHandler(ts.URL, "token", "!room:example.org")
	changes := []monitor.Change{
		matrixChange(monitor.ChangeAdded, 443, "nginx"),
	}
	if err := h.Handle(changes); err == nil {
		t.Error("expected error on non-2xx response")
	}
}

func TestFormatMatrixMsg_Added(t *testing.T) {
	c := matrixChange(monitor.ChangeAdded, 22, "sshd")
	msg := formatMatrixMsg(c)
	if !strings.Contains(msg, "added") {
		t.Errorf("expected 'added' in message, got %q", msg)
	}
	if !strings.Contains(msg, "22") {
		t.Errorf("expected port 22 in message, got %q", msg)
	}
}

func TestFormatMatrixMsg_FallsBackToIP(t *testing.T) {
	c := monitor.Change{
		Kind: monitor.ChangeRemoved,
		Binding: ports.Binding{
			Proto: "udp",
			Port:  53,
			Addr:  "192.168.1.1",
		},
	}
	msg := formatMatrixMsg(c)
	if !strings.Contains(msg, "192.168.1.1") {
		t.Errorf("expected IP in message, got %q", msg)
	}
	if !strings.Contains(msg, "removed") {
		t.Errorf("expected 'removed' in message, got %q", msg)
	}
}
