package alert

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/netwatch/portwatch/internal/monitor"
	"github.com/netwatch/portwatch/internal/ports"
)

// fakeProducer records published messages for inspection in tests.
type fakeProducer struct {
	msgs   []fakeMsg
	failOn int // if > 0, fail after this many successful publishes
}

type fakeMsg struct {
	topic string
	key   []byte
	value []byte
}

func (f *fakeProducer) Publish(topic string, key, value []byte) error {
	if f.failOn > 0 && len(f.msgs) >= f.failOn {
		return errors.New("producer error")
	}
	f.msgs = append(f.msgs, fakeMsg{topic: topic, key: key, value: value})
	return nil
}

func (f *fakeProducer) Close() error { return nil }

func kafkaChange(kind monitor.ChangeKind, port uint16, proto string) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{Proto: proto, Addr: "127.0.0.1", Port: port, Process: "sshd", PID: 42},
	}
}

func TestKafkaHandler_EmptyChangesNoPublish(t *testing.T) {
	p := &fakeProducer{}
	h := NewKafkaHandler(p, "portwatch")
	if err := h.Handle(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.msgs) != 0 {
		t.Fatalf("expected 0 messages, got %d", len(p.msgs))
	}
}

func TestKafkaHandler_PublishesOnChange(t *testing.T) {
	p := &fakeProducer{}
	h := NewKafkaHandler(p, "portwatch")
	changes := []monitor.Change{
		kafkaChange(monitor.Added, 22, "tcp"),
		kafkaChange(monitor.Removed, 8080, "tcp"),
	}
	if err := h.Handle(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(p.msgs))
	}
	for _, m := range p.msgs {
		if m.topic != "portwatch" {
			t.Errorf("expected topic %q, got %q", "portwatch", m.topic)
		}
		var ev kafkaEvent
		if err := json.Unmarshal(m.value, &ev); err != nil {
			t.Fatalf("invalid JSON payload: %v", err)
		}
		if ev.Timestamp == "" {
			t.Error("expected non-empty timestamp")
		}
		if ev.Process != "sshd" {
			t.Errorf("expected process %q, got %q", "sshd", ev.Process)
		}
	}
}

func TestKafkaHandler_ErrorOnProducerFailure(t *testing.T) {
	p := &fakeProducer{failOn: 1}
	h := NewKafkaHandler(p, "portwatch")
	changes := []monitor.Change{
		kafkaChange(monitor.Added, 443, "tcp"),
		kafkaChange(monitor.Added, 80, "tcp"),
	}
	if err := h.Handle(changes); err == nil {
		t.Fatal("expected error from failing producer")
	}
}

func TestKafkaHandler_DrainIsNoop(t *testing.T) {
	h := NewKafkaHandler(&fakeProducer{}, "portwatch")
	if err := h.Drain(); err != nil {
		t.Fatalf("Drain returned unexpected error: %v", err)
	}
}
