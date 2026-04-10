// Package ports provides utilities for scanning and resolving network port bindings.
package ports

import (
	"fmt"
	"net"
	"strconv"
)

// ResolvedBinding extends a Binding with human-readable service name and hostname.
type ResolvedBinding struct {
	Binding
	ServiceName string
	Hostname    string
}

// Resolver performs reverse DNS and service name lookups for port bindings.
type Resolver struct {
	lookupAddr func(string) ([]string, error)
	lookupPort func(string, int) string
}

// NewResolver creates a Resolver using the standard net.
func NewResolver() *Resolver {
	return &Resolver{
		lookupAddr: net.LookupAddr,
		lookupPort: lookupServiceName,
	}
}

//es a Binding with reverse DNS and service name information.
func (r *Resolver) Resolve(b Binding) ResolvedBinding {
	rb := ResolvedBinding{Binding: b}

	hosts, err := r.lookupAddr(b.IP)
	if err == nil && len(hosts) > 0 {
		rb.Hostname = hosts[0]
	} else {
		rb.Hostname = b.IP
	}

	rb.ServiceName = r.lookupPort(b.Proto, b.Port)
	return rb
}

// ResolveAll enriches a slice of Bindings.
func (r *Resolver) ResolveAll(bindings []Binding) []ResolvedBinding {
	result := make([]ResolvedBinding, 0, len(bindings))
	for _, b := range bindings {
		result = append(result, r.Resolve(b))
	}
	return result
}

// lookupServiceName returns a well-known service name for a port, or a numeric string.
func lookupServiceName(proto string, port int) string {
	name, err := net.LookupPort(proto, strconv.Itoa(port))
	if err != nil || name == 0 {
		return fmt.Sprintf("%d", port)
	}
	// net.LookupPort returns the port number; use getservbyport-style lookup via IANA names.
	// Fall back to well-known table for common ports.
	if svc, ok := wellKnown[port]; ok {
		return svc
	}
	return fmt.Sprintf("%d", port)
}

// wellKnown maps common port numbers to IANA service names.
var wellKnown = map[int]string{
	21:   "ftp",
	22:   "ssh",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	3306: "mysql",
	5432: "postgresql",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}
