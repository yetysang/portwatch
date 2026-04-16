package config

import "fmt"

// SyslogConfig holds configuration for the syslog alert handler.
type SyslogConfig struct {
	Enabled  bool   `toml:"enabled"`
	Network  string `toml:"network"`  // "tcp", "udp", or "" for local
	Addr     string `toml:"addr"`     // host:port, empty for local
	Tag      string `toml:"tag"`
	Facility string `toml:"facility"`
}

// DefaultSyslogConfig returns a SyslogConfig with sensible defaults.
func DefaultSyslogConfig() SyslogConfig {
	return SyslogConfig{
		Enabled:  false,
		Network:  "",
		Addr:     "",
		Tag:      "portwatch",
		Facility: "daemon",
	}
}

// Validate checks that the SyslogConfig is well-formed.
func (c SyslogConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Tag == "" {
		return fmt.Errorf("syslog: tag must not be empty")
	}
	validFacilities := map[string]bool{
		"kern": true, "user": true, "mail": true, "daemon": true,
		"auth": true, "syslog": true, "lpr": true, "news": true,
		"local0": true, "local1": true, "local2": true, "local3": true,
		"local4": true, "local5": true, "local6": true, "local7": true,
	}
	if !validFacilities[c.Facility] {
		return fmt.Errorf("syslog: unknown facility %q", c.Facility)
	}
	if (c.Network == "") != (c.Addr == "") {
		return fmt.Errorf("syslog: network and addr must both be set or both be empty")
	}
	return nil
}
