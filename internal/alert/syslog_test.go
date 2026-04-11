package alert

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

func syslogChange(kind monitor.ChangeKind, port int) monitor.Change {
	return monitor.Change{
		Kind: kind,
		Binding: ports.Binding{
			Proto:    "tcp",
			Host:     "127.0.0.1",
			Port:     port,
			Hostname: "localhost",
			Service:  "http",
			PID:      1234,
			Process:  "nginx",
		},
	}
}

func TestFormatSyslogMsg_Added(t *testing.T) {
	c := syslogChange(monitor.Added, 80)
	msg := formatSyslogMsg(c)
	if !strings.Contains(msg, "[added]") {
		t.Errorf("expected '[added]' in message, got: %s", msg)
	}
	if !strings.Contains(msg, "TCP") {
		t.Errorf("expected 'TCP' in message, got: %s", msg)
	}
	if !strings.Contains(msg, "localhost:http") {
		t.Errorf("expected 'localhost:http' in message, got: %s", msg)
	}
	if !strings.Contains(msg, "pid=1234") {
		t.Errorf("expected 'pid=1234' in message, got: %s", msg)
	}
	if !strings.Contains(msg, "process=nginx") {
		t.Errorf("expected 'process=nginx' in message, got: %s", msg)
	}
}

func TestFormatSyslogMsg_FallsBackToIP(t *testing.T) {
	c := monitor.Change{
		Kind: monitor.Removed,
		Binding: ports.Binding{
			Proto: "udp",
			Host:  "0.0.0.0",
			Port:  5353,
		},
	}
	msg := formatSyslogMsg(c)
	if !strings.Contains(msg, "0.0.0.0:5353") {
		t.Errorf("expected '0.0.0.0:5353' in message, got: %s", msg)
	}
	if !strings.Contains(msg, "[removed]") {
		t.Errorf("expected '[removed]' in message, got: %s", msg)
	}
}

func TestFormatSyslogMsg_NoPIDOrProcess(t *testing.T) {
	c := monitor.Change{
		Kind: monitor.Added,
		Binding: ports.Binding{
			Proto:   "tcp",
			Host:    "0.0.0.0",
			Port:    443,
			Service: "https",
		},
	}
	msg := formatSyslogMsg(c)
	if strings.Contains(msg, "pid=") {
		t.Errorf("did not expect 'pid=' when PID is 0, got: %s", msg)
	}
	if strings.Contains(msg, "process=") {
		t.Errorf("did not expect 'process=' when Process is empty, got: %s", msg)
	}
}

func TestFormatSyslogMsg_ProtoUpperCase(t *testing.T) {
	c := monitor.Change{
		Kind: monitor.Added,
		Binding: ports.Binding{Proto: "tcp", Host: "::1", Port: 22},
	}
	msg := formatSyslogMsg(c)
	if !strings.Contains(msg, "TCP") {
		t.Errorf("expected uppercase 'TCP', got: %s", msg)
	}
}
