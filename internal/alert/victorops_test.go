package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"portwatch/internal/monitor"
)

func victoropsChange(kind monitor.ChangeKind, port uint16, proto, ip string) monitor.Change {
	return monitor.Change{
		Kind:    kind,
		Binding: bindingFor(port, proto, ip),
	}
}

func TestVictorOpsHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	h := NewVictorOpsHandler(ts.URL)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request for empty changes")
	}
}

func TestVictorOpsHandler_PostsOnChange(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewVictorOpsHandler(ts.URL)
	changes := []monitor.Change{
		victoropsChange(monitor.ChangeAdded, 8080, "tcp", "0.0.0.0"),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["message_type"] != "CRITICAL" {
		t.Errorf("expected CRITICAL message_type, got %v", received["message_type"])
	}
}

func TestVictorOpsHandler_RecoveryOnRemoved(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewVictorOpsHandler(ts.URL)
	changes := []monitor.Change{
		victoropsChange(monitor.ChangeRemoved, 22, "tcp", "0.0.0.0"),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["message_type"] != "RECOVERY" {
		t.Errorf("expected RECOVERY message_type, got %v", received["message_type"])
	}
}

func TestVictorOpsHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewVictorOpsHandler(ts.URL)
	changes := []monitor.Change{
		victoropsChange(monitor.ChangeAdded, 9090, "tcp", "127.0.0.1"),
	}
	if err := h.Handle(changes); err == nil {
		t.Error("expected error on non-2xx response")
	}
}

func TestFormatVictorOpsMsg_Added(t *testing.T) {
	c := victoropsChange(monitor.ChangeAdded, 443, "tcp", "0.0.0.0")
	msg := formatVictorOpsMsg(c)
	if msg == "" {
		t.Error("expected non-empty message")
	}
}

func TestVictorOpsHandler_DrainIsNoop(t *testing.T) {
	h := NewVictorOpsHandler("http://localhost")
	if err := h.Drain(); err != nil {
		t.Fatalf("unexpected error from Drain: %v", err)
	}
}
