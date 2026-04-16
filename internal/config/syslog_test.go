package config

import "testing"

func TestDefaultSyslogConfig_Values(t *testing.T) {
	c := DefaultSyslogConfig()
	if c.Enabled {
		t.Error("expected Enabled=false")
	}
	if c.Tag != "portwatch" {
		t.Errorf("unexpected Tag: %q", c.Tag)
	}
	if c.Facility != "daemon" {
		t.Errorf("unexpected Facility: %q", c.Facility)
	}
	if c.Network != "" || c.Addr != "" {
		t.Error("expected empty Network and Addr")
	}
}

func TestSyslogConfig_ValidateDisabledSkipsChecks(t *testing.T) {
	c := SyslogConfig{Enabled: false, Tag: "", Facility: ""}
	if err := c.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSyslogConfig_ValidateEnabledRequiresTag(t *testing.T) {
	c := DefaultSyslogConfig()
	c.Enabled = true
	c.Tag = ""
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for empty tag")
	}
}

func TestSyslogConfig_ValidateInvalidFacility(t *testing.T) {
	c := DefaultSyslogConfig()
	c.Enabled = true
	c.Facility = "bogus"
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for invalid facility")
	}
}

func TestSyslogConfig_ValidateNetworkWithoutAddr(t *testing.T) {
	c := DefaultSyslogConfig()
	c.Enabled = true
	c.Network = "tcp"
	c.Addr = ""
	if err := c.Validate(); err == nil {
		t.Fatal("expected error when network set but addr empty")
	}
}

func TestSyslogConfig_ValidateRemoteValid(t *testing.T) {
	c := DefaultSyslogConfig()
	c.Enabled = true
	c.Network = "udp"
	c.Addr = "logs.example.com:514"
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSyslogConfig_ValidateLocalValid(t *testing.T) {
	c := DefaultSyslogConfig()
	c.Enabled = true
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
