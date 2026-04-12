package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"portwatch/internal/monitor"
)

func opsgenieChange(kind monitor.ChangeKind, port uint16, proto, ip string) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: bindingFor(port, proto, ip),
	}
}

func TestOpsGenieHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	h := NewOpsGenieHandler(ts.URL, "test-key")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request for empty changes")
	}
}

func TestOpsGenieHandler_PostsOnChange(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	h := NewOpsGenieHandler(ts.URL, "test-key")
	changes := []monitor.Change{
		opsgenieChange(monitor.ChangeAdded, 8080, "tcp", "0.0.0.0"),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["message"] == nil {
		t.Error("expected message field in payload")
	}
}

func TestOpsGenieHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	h := NewOpsGenieHandler(ts.URL, "bad-key")
	changes := []monitor.Change{
		opsgenieChange(monitor.ChangeAdded, 9090, "tcp", "127.0.0.1"),
	}
	if err := h.Handle(changes); err == nil {
		t.Error("expected error on non-2xx response")
	}
}

func TestFormatOpsGenieMsg_Added(t *testing.T) {
	c := opsgenieChange(monitor.ChangeAdded, 443, "tcp", "0.0.0.0")
	msg := formatOpsGenieMsg(c)
	if msg == "" {
		t.Error("expected non-empty message")
	}
}

func TestFormatOpsGenieMsg_Removed(t *testing.T) {
	c := opsgenieChange(monitor.ChangeRemoved, 22, "tcp", "0.0.0.0")
	msg := formatOpsGenieMsg(c)
	if msg == "" {
		t.Error("expected non-empty message for removed change")
	}
}

func TestOpsGenieHandler_DrainIsNoop(t *testing.T) {
	h := NewOpsGenieHandler("http://localhost", "key")
	if err := h.Drain(); err != nil {
		t.Fatalf("unexpected error from Drain: %v", err)
	}
}
