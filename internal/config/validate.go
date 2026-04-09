package config

import (
	"errors"
	"fmt"
	"time"
)

// ValidationError holds a list of validation failures for a Config.
type ValidationError struct {
	Errs []string
}

func (v *ValidationError) Error() string {
	if len(v.Errs) == 1 {
		return fmt.Sprintf("config validation error: %s", v.Errs[0])
	}
	return fmt.Sprintf("config validation errors (%d): %v", len(v.Errs), v.Errs)
}

// IsValidationError reports whether err is a *ValidationError.
func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}

// Validate checks that c contains sensible values and returns a
// *ValidationError describing every problem found, or nil on success.
func Validate(c *Config) error {
	var errs []string

	if c.Interval < 500*time.Millisecond {
		errs = append(errs, fmt.Sprintf("interval %v is below minimum 500ms", c.Interval))
	}
	if c.Interval > 24*time.Hour {
		errs = append(errs, fmt.Sprintf("interval %v exceeds maximum 24h", c.Interval))
	}

	if c.LogLevel != "info" && c.LogLevel != "warn" && c.LogLevel != "error" && c.LogLevel != "debug" {
		errs = append(errs, fmt.Sprintf("unknown log level %q (must be debug|info|warn|error)", c.LogLevel))
	}

	for _, p := range c.IgnorePorts {
		if p < 1 || p > 65535 {
			errs = append(errs, fmt.Sprintf("ignore port %d out of valid range 1-65535", p))
		}
	}

	if len(errs) > 0 {
		return &ValidationError{Errs: errs}
	}
	return nil
}
