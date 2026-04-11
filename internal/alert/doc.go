// Package alert implements change notification handlers for portwatch.
//
// The primary type is Handler, which formats and buffers alert lines for
// consumption by the CLI layer. Each alert line is prefixed with a severity
// level (INFO or WARN) derived from the nature of the port change.
//
// Severity levels
//
// INFO is used for ports that have been closed or released. WARN is used for
// newly opened ports, since unexpected listeners may indicate a security
// concern or misconfiguration.
//
// Rate-limiting
//
// RateLimitedHandler wraps a Handler with a ports.RateLimiter so that
// repeated alerts for the same port/protocol pair are suppressed during a
// configurable cooldown window. This prevents log spam when a process
// repeatedly opens and closes the same port in quick succession.
//
// Usage:
//
//	rl := ports.NewRateLimiter(30 * time.Second)
//	h := alert.NewRateLimitedHandler(alert.NewHandler(cfg), rl)
//	h.Handle(changes)
//	for _, line := range h.Drain() {
//		fmt.Println(line)
//	}
package alert
