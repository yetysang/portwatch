package alert

import (
	"log/syslog"
	"testing"
)

func TestResolveFacility_Daemon(t *testing.T) {
	p, err := resolveFacility("daemon")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != syslog.LOG_DAEMON {
		t.Errorf("expected LOG_DAEMON, got %v", p)
	}
}

func TestResolveFacility_Local3(t *testing.T) {
	p, err := resolveFacility("local3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != syslog.LOG_LOCAL3 {
		t.Errorf("expected LOG_LOCAL3, got %v", p)
	}
}

func TestResolveFacility_Unknown(t *testing.T) {
	_, err := resolveFacility("bogus")
	if err == nil {
		t.Fatal("expected error for unknown facility")
	}
}

func TestResolveFacility_AllKnown(t *testing.T) {
	known := []string{
		"kern", "user", "mail", "daemon", "auth", "syslog",
		"lpr", "news",
		"local0", "local1", "local2", "local3",
		"local4", "local5", "local6", "local7",
	}
	for _, name := range known {
		_, err := resolveFacility(name)
		if err != nil {
			t.Errorf("unexpected error for facility %q: %v", name, err)
		}
	}
}
