package alert

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/iamcalledned/portwatch/internal/monitor"
	"github.com/iamcalledned/portwatch/internal/ports"
)

// stubAMQPPublisher records published messages for inspection.
type stubAMQPPublisher struct {
	published [][]byte
	failOn    int // if > 0, fail on the nth Publish call (1-indexed)
	calls     int
}

func (s *stubAMQPPublisher) Publish(_, _ string, body []byte) error {
	s.calls++
	if s.failOn > 0 && s.calls == s.failOn {
		return errors.New("broker unavailable")
	}
	s.published = append(s.published, body)
	return nil
}
func (s *stubAMQPPublisher) Close() error { return nil }

func amqpChange(kind monitor.ChangeKind, port uint16) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Proto:    "tcp",
			Addr:     "127.0.0.1",
			Hostname: "localhost",
			Port:     port,
			Process:  "nginx",
			PID:      1234,
		},
	}
}

func TestAMQPHandler_EmptyChangesNoPublish(t *testing.T) {
	pub := &stubAMQPPublisher{}
	h := NewAMQPHandler(pub, "portwatch", "events")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pub.calls != 0 {
		t.Errorf("expected 0 publish calls, got %d", pub.calls)
	}
}

func TestAMQPHandler_PublishesOnChange(t *testing.T) {
	pub := &stubAMQPPublisher{}
	h := NewAMQPHandler(pub, "portwatch", "events")
	changes := []monitor.Change{
		amqpChange(monitor.ChangeKindAdded, 8080),
		amqpChange(monitor.ChangeKindRemoved, 9090),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pub.published) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(pub.published))
	}
}

func TestAMQPHandler_MessageFields(t *testing.T) {
	pub := &stubAMQPPublisher{}
	h := NewAMQPHandler(pub, "portwatch", "events")
	if err := h.Handle([]monitor.Change{amqpChange(monitor.ChangeKindAdded, 443)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var msg amqpMessage
	if err := json.Unmarshal(pub.published[0], &msg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if msg.Port != 443 {
		t.Errorf("expected port 443, got %d", msg.Port)
	}
	if msg.Action != "added" {
		t.Errorf("expected action 'added', got %q", msg.Action)
	}
	if msg.Addr != "localhost" {
		t.Errorf("expected addr 'localhost', got %q", msg.Addr)
	}
	if msg.Process != "nginx" {
		t.Errorf("expected process 'nginx', got %q", msg.Process)
	}
}

func TestAMQPHandler_ErrorOnPublishFailure(t *testing.T) {
	pub := &stubAMQPPublisher{failOn: 1}
	h := NewAMQPHandler(pub, "portwatch", "events")
	err := h.Handle([]monitor.Change{amqpChange(monitor.ChangeKindAdded, 80)})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAMQPHandler_DrainIsNoop(t *testing.T) {
	pub := &stubAMQPPublisher{}
	h := NewAMQPHandler(pub, "portwatch", "events")
	if err := h.Drain(); err != nil {
		t.Errorf("Drain returned unexpected error: %v", err)
	}
}
