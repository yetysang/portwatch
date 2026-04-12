package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func pdChange(kind monitor.ChangeKind) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Addr:     "127.0.0.1",
			Hostname: "localhost",
			Port:     "9200",
			Proto:    "tcp",
			PID:      42,
			Process:  "elasticsearch",
		},
	}
}

func TestPagerDutyHandler_EmptyChangesNoRequest(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()

	h := NewPagerDutyHandler("key123", "error")
	h.client = ts.Client()
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected no HTTP request for empty changes")
	}
}

func TestPagerDutyHandler_PostsOnChange(t *testing.T) {
	var captured pdPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	h := NewPagerDutyHandler("routekey", "critical")
	h.client = ts.Client()
	// override URL via a wrapper — we test the real path by pointing client at ts
	// We need to redirect the URL; swap out via a custom RoundTripper.
	h.client = &http.Client{Transport: redirectTransport(ts.URL)}

	if err := h.Handle([]monitor.Change{pdChange(monitor.Added)}); err != nil {
		t.Fatalf("Handle error: %v", err)
	}
	if captured.RoutingKey != "routekey" {
		t.Errorf("routing key = %q, want %q", captured.RoutingKey, "routekey")
	}
	if captured.Payload.Severity != "critical" {
		t.Errorf("severity = %q, want critical", captured.Payload.Severity)
	}
	if captured.EventAction != "trigger" {
		t.Errorf("event_action = %q, want trigger", captured.EventAction)
	}
}

func TestPagerDutyHandler_ErrorOnNonSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	h := NewPagerDutyHandler("key", "")
	h.client = &http.Client{Transport: redirectTransport(ts.URL)}
	err := h.Handle([]monitor.Change{pdChange(monitor.Removed)})
	if err == nil {
		t.Fatal("expected error on 400 response")
	}
}

func TestFormatPagerDutyMsg_Added(t *testing.T) {
	msg := formatPagerDutyMsg(pdChange(monitor.Added))
	for _, want := range []string{"9200", "tcp", "added", "localhost", "elasticsearch"} {
		if !containsStr(msg, want) {
			t.Errorf("message %q missing %q", msg, want)
		}
	}
}

func TestFormatPagerDutyMsg_FallsBackToIP(t *testing.T) {
	c := pdChange(monitor.Removed)
	c.Binding.Hostname = ""
	msg := formatPagerDutyMsg(c)
	if !containsStr(msg, "127.0.0.1") {
		t.Errorf("expected IP fallback in %q", msg)
	}
}

// redirectTransport rewrites every request URL to the given base.
type redirectTransport string

func (rt redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req.Clone(req.Context())
	req2.URL.Scheme = "http"
	req2.URL.Host = string(rt)[len("http://"):]
	return http.DefaultTransport.RoundTrip(req2)
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && func() bool {
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	}())
}
