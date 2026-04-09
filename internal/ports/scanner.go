package ports

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Binding represents a single port binding observed on the system.
type Binding struct {
	Protocol string
	LocalAddr string
	Port     uint16
	PID      int
	State    string
}

// Scanner reads current port bindings from the OS.
type Scanner struct{}

// NewScanner creates a new Scanner instance.
func NewScanner() *Scanner {
	return &Scanner{}
}

// Scan returns all current TCP/UDP port bindings by reading /proc/net.
func (s *Scanner) Scan() ([]Binding, error) {
	var bindings []Binding

	for _, proto := range []string{"tcp", "tcp6", "udp", "udp6"} {
		path := fmt.Sprintf("/proc/net/%s", proto)
		b, err := s.parseNetFile(path, proto)
		if err != nil {
			// Non-fatal: file may not exist on all systems
			continue
		}
		bindings = append(bindings, b...)
	}

	return bindings, nil
}

func (s *Scanner) parseNetFile(path, proto string) ([]Binding, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var bindings []Binding
	scanner := bufio.NewScanner(f)

	// Skip header line
	scanner.Scan()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		localAddr, port, err := parseHexAddr(fields[1])
		if err != nil {
			continue
		}

		state := fields[3]

		bindings = append(bindings, Binding{
			Protocol:  proto,
			LocalAddr: localAddr,
			Port:      port,
			State:     state,
		})
	}

	return bindings, scanner.Err()
}

// parseHexAddr parses a hex-encoded "addr:port" field from /proc/net files.
func parseHexAddr(hexAddr string) (string, uint16, error) {
	parts := strings.Split(hexAddr, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid address format: %s", hexAddr)
	}

	portVal, err := strconv.ParseUint(parts[1], 16, 16)
	if err != nil {
		return "", 0, err
	}

	return parts[0], uint16(portVal), nil
}
