// Package ports provides primitives for discovering, filtering, and enriching
// network port binding information on the local system.
//
// # Scanner
//
// NewScanner returns a Scanner that reads active TCP/UDP bindings from the
// /proc/net filesystem. Call Scan() to obtain a snapshot of current bindings
// as a map keyed by "proto:ip:port".
//
// # Filter
//
// NewFilter wraps an IgnoreSet and removes user-configured ports from a
// binding map before it is handed to the monitor.
//
// # Resolver
//
// NewResolver enriches Binding values with reverse-DNS hostnames and IANA
// service names. Resolution is best-effort; failures fall back to the raw IP
// address or numeric port string.
//
// # Proc
//
// LookupProc resolves the owning process (PID + command name) for a socket
// inode by walking /proc/*/fd symlinks.
package ports
