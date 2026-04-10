// Package ports provides scanning and filtering of active network port
// bindings on the local host.
//
// # Scanner
//
// NewScanner returns a Scanner that reads from /proc/net/tcp and
// /proc/net/tcp6 (on Linux) to enumerate listening TCP sockets. Each
// socket is returned as a Binding value containing the protocol, address,
// port and owning PID.
//
// # Filter
//
// NewFilter wraps a config.IgnoreSet and exposes Apply / ApplyToMap helpers
// that strip any Binding whose port appears in the ignore list. This lets
// callers exclude well-known or expected ports (e.g. 22, 80) before passing
// results to the monitor or alert subsystems.
//
// Typical usage:
//
//	scanner := ports.NewScanner()
//	filter  := ports.NewFilter(ignoreSet)
//
//	bindings, err := scanner.Scan()
//	if err != nil { ... }
//	bindings = filter.Apply(bindings)
package ports
