package alert

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func telegramChange(kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Addr:     "127.0.0.1",
			Hostname: "localhost",
			Port:     9090,
			Proto:    "tcp",
			Process:  "myapp",
			PID:      1234,
		},
	}
}

func TestTelegramHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	h := NewTelegramHandler("token", "chat123", ts.Client())
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected no HTTP request for empty changes")
	}
}

func TestTelegramHandler_PostsOnChange(t *testing.T) {
	var received telegramPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Override the API URL by using a custom client that redirects to the test server.
	h := NewTelegramHandler("mytoken", "chat42", ts.Client())
	// Patch the URL by wrapping the client transport.
	h.client = &http.Client{
		Transport: rewriteTransport{base: ts.Client().Transport, target: ts.URL},
	}

	changes := []monitor.Change{telegramChange(monitor.Added)}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.ChatID != "chat42" {
		t.Errorf("expected chat_id 'chat42', got %q", received.ChatID)
	}
	if !strings.Contains(received.Text, "9090") {
		t.Errorf("expected port 9090 in message, got %q", received.Text)
	}
	if !strings.Contains(received.Text, "added") {
		t.Errorf("expected 'added' in message, got %q", received.Text)
	}
}

func TestTelegramHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	h := NewTelegramHandler("badtoken", "chat1", ts.Client())
	h.client = &http.Client{
		Transport: rewriteTransport{base: ts.Client().Transport, target: ts.URL},
	}

	changes := []monitor.Change{telegramChange(monitor.Removed)}
	err := h.Handle(changes)
	if err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}

func TestFormatTelegramMsg_Added(t *testing.T) {
	c := telegramChange(monitor.Added)
	msg := formatTelegramMsg(c)
	if !strings.Contains(msg, "added") {
		t.Errorf("expected 'added' in message: %q", msg)
	}
	if !strings.Contains(msg, "9090") {
		t.Errorf("expected port in message: %q", msg)
	}
	if !strings.Contains(msg, "myapp") {
		t.Errorf("expected process name in message: %q", msg)
	}
}

func TestFormatTelegramMsg_FallsBackToIP(t *testing.T) {
	c := monitor.Change{
		Kind: monitor.Added,
		Binding: ports.Binding{Addr: "10.0.0.1", Port: 8080, Proto: "udp"},
	}
	msg := formatTelegramMsg(c)
	if !strings.Contains(msg, "10.0.0.1") {
		t.Errorf("expected IP fallback in message: %q", msg)
	}
}

// rewriteTransport redirects all requests to a fixed base URL (for testing).
type rewriteTransport struct {
	base   http.RoundTripper
	target string
}

func (rt rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req.Clone(req.Context())
	req2.URL.Scheme = "http"
	req2.URL.Host = strings.TrimPrefix(rt.target, "http://")
	if rt.base != nil {
		return rt.base.RoundTrip(req2)
	}
	return http.DefaultTransport.RoundTrip(req2)
}
