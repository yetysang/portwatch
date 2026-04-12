package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func splunkChange(kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Addr:     "127.0.0.1",
			Hostname: "localhost",
			Port:     9200,
			Proto:    "tcp",
			PID:      42,
			Process:  "elasticsearch",
			Service:  "elasticsearch",
		},
	}
}

func TestSplunkHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer srv.Close()

	h := NewSplunkHandler(srv.URL, "mytoken")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected no HTTP request for empty changes")
	}
}

func TestSplunkHandler_PostsOnChange(t *testing.T) {
	var gotAuth string
	var body []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		body, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h := NewSplunkHandler(srv.URL, "secret")
	if err := h.Handle([]monitor.Change{splunkChange(monitor.ChangeAdded)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAuth != "Splunk secret" {
		t.Errorf("expected Authorization header 'Splunk secret', got %q", gotAuth)
	}
	var ev map[string]any
	if err := json.Unmarshal(body, &ev); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if ev["source"] != "portwatch" {
		t.Errorf("expected source=portwatch, got %v", ev["source"])
	}
}

func TestSplunkHandler_ErrorOnNonSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	h := NewSplunkHandler(srv.URL, "bad")
	err := h.Handle([]monitor.Change{splunkChange(monitor.ChangeAdded)})
	if err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestFormatSplunkEvent_FallsBackToIP(t *testing.T) {
	c := monitor.Change{
		Kind: monitor.ChangeAdded,
		Binding: ports.Binding{
			Addr:  "10.0.0.1",
			Port:  8080,
			Proto: "tcp",
		},
	}
	ev := formatSplunkEvent(c)
	if ev["host"] != "10.0.0.1" {
		t.Errorf("expected host fallback to IP, got %v", ev["host"])
	}
}
