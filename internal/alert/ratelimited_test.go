package alert

import (
	"strings"
	"testing"
	"time"

	"github.com/example/portwatch/internal/monitor"
	"github.com/example/portwatch/internal/ports"
)

func makeRLHandler(cooldown time.Duration) (*RateLimitedHandler, *ports.RateLimiter) {
	cfg := defaultTestConfig()
	h := NewHandler(cfg)
	rl := ports.NewRateLimiter(cooldown)
	return NewRateLimitedHandler(h, rl), rl
}

func defaultTestConfig() interface{ LogLevel() string } {
	return nil
}

func rlChange(kind monitor.ChangeKind, proto string, port int) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{Proto: proto, Port: port, Addr: "127.0.0.1"},
	}
}

func TestRateLimitedHandler_FirstChangeForwarded(t *testing.T) {
	h, _ := makeRLHandler(10 * time.Second)
	h.Handle([]monitor.Change{rlChange(monitor.Added, "tcp", 8080)})
	out := h.Drain()
	if len(out) != 1 {
		t.Fatalf("expected 1 output line, got %d", len(out))
	}
}

func TestRateLimitedHandler_DuplicateSuppressed(t *testing.T) {
	h, _ := makeRLHandler(10 * time.Second)
	c := rlChange(monitor.Added, "tcp", 8080)
	h.Handle([]monitor.Change{c})
	h.Drain()
	h.Handle([]monitor.Change{c})
	out := h.Drain()
	if len(out) != 0 {
		t.Fatalf("expected duplicate to be suppressed, got %d lines", len(out))
	}
}

func TestRateLimitedHandler_DifferentPortsForwarded(t *testing.T) {
	h, _ := makeRLHandler(10 * time.Second)
	h.Handle([]monitor.Change{
		rlChange(monitor.Added, "tcp", 8080),
		rlChange(monitor.Added, "tcp", 9090),
	})
	out := h.Drain()
	if len(out) != 2 {
		t.Fatalf("expected 2 output lines, got %d", len(out))
	}
}

func TestRateLimitedHandler_ResetAllowsReplay(t *testing.T) {
	h, rl := makeRLHandler(10 * time.Second)
	c := rlChange(monitor.Added, "tcp", 8080)
	h.Handle([]monitor.Change{c})
	h.Drain()
	rl.Reset()
	h.Handle([]monitor.Change{c})
	out := h.Drain()
	if len(out) != 1 {
		t.Fatalf("expected 1 line after reset, got %d", len(out))
	}
}

func TestChangeKey_Format(t *testing.T) {
	c := rlChange(monitor.Added, "tcp", 443)
	key := changeKey(c)
	if !strings.Contains(key, "tcp") || !strings.Contains(key, "443") {
		t.Fatalf("unexpected key format: %s", key)
	}
}
