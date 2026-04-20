package config

// SignalConfig holds configuration for OS signal handling behaviour.
type SignalConfig struct {
	// GracePeriodSeconds is how long the daemon waits for in-flight work to
	// finish after receiving SIGTERM or SIGINT before forcefully exiting.
	GracePeriodSeconds int `toml:"grace_period_seconds" json:"grace_period_seconds"`

	// ReloadOnHUP controls whether SIGHUP triggers a config reload.
	ReloadOnHUP bool `toml:"reload_on_hup" json:"reload_on_hup"`
}

// DefaultSignalConfig returns a SignalConfig with sensible defaults.
func DefaultSignalConfig() SignalConfig {
	return SignalConfig{
		GracePeriodSeconds: 5,
		ReloadOnHUP:        true,
	}
}

// Validate checks that SignalConfig values are within acceptable ranges.
func (s SignalConfig) Validate() error {
	if s.GracePeriodSeconds < 0 {
		return ValidationError{Field: "grace_period_seconds", Message: "must be non-negative"}
	}
	if s.GracePeriodSeconds > 120 {
		return ValidationError{Field: "grace_period_seconds", Message: "must not exceed 120"}
	}
	return nil
}
