package alert

import (
	"encoding/json"
	"errors"
	"testing"

	"portwatch/internal/monitor"
	"portwatch/internal/ports"
)

// stubNatsConn implements natsPublisher for testing.
type stubNatsConn struct {
	published [][]byte
	subjects  []string
	errPublish error
	errDrain   error
}

func (s *stubNatsConn) Publish(subject string, data []byte) error {
	if s.errPublish != nil {
		return s.errPublish
	}
	s.subjects = append(s.subjects, subject)
	s.published = append(s.published, data)
	return nil
}

func (s *stubNatsConn) Drain() error { return s.errDrain }

func natsChange() monitor.Change {
	return monitor.Change{
		Kind: monitor.Added,
		Binding: ports.Binding{
			Proto:   "tcp",
			Addr:    "0.0.0.0",
			Port:    9090,
			Process: "myapp",
			PID:     1234,
		},
	}
}

func TestNatsHandler_EmptyChangesNoPublish(t *testing.T) {
	conn := &stubNatsConn{}
	h := &NatsHandler{conn: conn, subject: "portwatch.events"}
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conn.published) != 0 {
		t.Errorf("expected no publishes, got %d", len(conn.published))
	}
}

func TestNatsHandler_PublishesOnChange(t *testing.T) {
	conn := &stubNatsConn{}
	h := &NatsHandler{conn: conn, subject: "portwatch.events"}
	c := natsChange()
	if err := h.Handle([]monitor.Change{c}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conn.published) != 1 {
		t.Fatalf("expected 1 publish, got %d", len(conn.published))
	}
	if conn.subjects[0] != "portwatch.events" {
		t.Errorf("unexpected subject: %s", conn.subjects[0])
	}
	var payload NatsPayload
	if err := json.Unmarshal(conn.published[0], &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload.Kind != "added" {
		t.Errorf("expected kind=added, got %s", payload.Kind)
	}
	if payload.Port != 9090 {
		t.Errorf("expected port=9090, got %d", payload.Port)
	}
	if payload.Process != "myapp" {
		t.Errorf("expected process=myapp, got %s", payload.Process)
	}
}

func TestNatsHandler_ErrorOnPublishFailure(t *testing.T) {
	conn := &stubNatsConn{errPublish: errors.New("connection refused")}
	h := &NatsHandler{conn: conn, subject: "portwatch.events"}
	err := h.Handle([]monitor.Change{natsChange()})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNatsHandler_DrainDelegates(t *testing.T) {
	conn := &stubNatsConn{errDrain: errors.New("drain failed")}
	h := &NatsHandler{conn: conn, subject: "portwatch.events"}
	if err := h.Drain(); err == nil {
		t.Error("expected drain error, got nil")
	}
}
