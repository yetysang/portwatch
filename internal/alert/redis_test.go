package alert

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

type stubRedis struct {
	calls []map[string]interface{}
	err   error
}

func (s *stubRedis) XAdd(_ context.Context, _, _ string, values map[string]interface{}) error {
	if s.err != nil {
		return s.err
	}
	s.calls = append(s.calls, values)
	return nil
}

func (s *stubRedis) Close() error { return nil }

func redisChange(kind monitor.ChangeKind, port uint16) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{Port: port, Proto: "tcp", Addr: "0.0.0.0", Hostname: "host1"},
	}
}

func TestRedisHandler_EmptyChangesNoPublish(t *testing.T) {
	stub := &stubRedis{}
	h := NewRedisHandler(stub, "portwatch")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stub.calls) != 0 {
		t.Errorf("expected 0 calls, got %d", len(stub.calls))
	}
}

func TestRedisHandler_PublishesOnChange(t *testing.T) {
	stub := &stubRedis{}
	h := NewRedisHandler(stub, "portwatch")
	changes := []monitor.Change{redisChange(monitor.Added, 8080)}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stub.calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(stub.calls))
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(stub.calls[0]["event"].(string)), &payload); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}
	if payload["kind"] != "added" {
		t.Errorf("expected kind=added, got %v", payload["kind"])
	}
}

func TestRedisHandler_ErrorOnFailure(t *testing.T) {
	stub := &stubRedis{err: context.DeadlineExceeded}
	h := NewRedisHandler(stub, "portwatch")
	changes := []monitor.Change{redisChange(monitor.Added, 9090)}
	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFormatRedisMsg_Fields(t *testing.T) {
	c := redisChange(monitor.Removed, 443)
	msg := formatRedisMsg(c)
	if _, ok := msg["event"]; !ok {
		t.Error("expected 'event' field")
	}
	if _, ok := msg["ts"]; !ok {
		t.Error("expected 'ts' field")
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(msg["event"].(string)), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload["kind"] != "removed" {
		t.Errorf("expected kind=removed, got %v", payload["kind"])
	}
}
