// Package ports provides types and utilities for scanning active TCP/UDP port
// bindings on the local system.
//
// # Scanner
//
// NewScanner returns a Scanner that reads from /proc/net/tcp and /proc/net/tcp6
// to enumerate listening ports. Each result is returned as a Binding value
// containing the protocol, address, and port.
//
// # Filter
//
// NewFilter wraps an IgnoreSet and can be used to exclude known or expected
// bindings from scan results before they are passed to the monitor.
//
// # Proc
//
// LookupProc resolves a PID to a human-readable process name by reading from
// the /proc filesystem. ParseInode extracts the inode field from a raw
// /proc/net entry, which can be used to correlate sockets with processes via
// /proc/<pid>/fd symlinks.
package ports
