package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func webhookChange(event monitor.ChangeType) monitor.Change {
	return monitor.Change{
		Type: event,
		Binding: ports.Binding{
			Proto:   "tcp",
			Addr:    "0.0.0.0",
			Port:    8080,
			Process: "myapp",
			PID:     1234,
		},
	}
}

func TestWebhookHandler_PostsOnChange(t *testing.T) {
	var received WebhookPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewWebhookHandler(ts.URL, 5*time.Second)
	if err := h.Handle([]monitor.Change{webhookChange(monitor.ChangeAdded)}); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}

	if received.Port != 8080 {
		t.Errorf("expected port 8080, got %d", received.Port)
	}
	if received.Event != string(monitor.ChangeAdded) {
		t.Errorf("expected event %q, got %q", monitor.ChangeAdded, received.Event)
	}
	if received.Process != "myapp" {
		t.Errorf("expected process myapp, got %q", received.Process)
	}
}

func TestWebhookHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewWebhookHandler(ts.URL, 5*time.Second)
	err := h.Handle([]monitor.Change{webhookChange(monitor.ChangeAdded)})
	if err == nil {
		t.Fatal("expected error for non-2xx response, got nil")
	}
}

func TestWebhookHandler_EmptyChangesNoRequest(t *testing.T) {
	requestCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewWebhookHandler(ts.URL, 5*time.Second)
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if requestCount != 0 {
		t.Errorf("expected 0 requests for empty changes, got %d", requestCount)
	}
}

func TestWebhookHandler_InvalidURLReturnsError(t *testing.T) {
	h := NewWebhookHandler("http://127.0.0.1:0/invalid", 500*time.Millisecond)
	err := h.Handle([]monitor.Change{webhookChange(monitor.ChangeRemoved)})
	if err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}
