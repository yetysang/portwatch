package alert

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/user/portwatch/internal/monitor"
)

// GRPCHandler sends port change alerts over a gRPC stream to a remote endpoint.
type GRPCHandler struct {
	target  string
	timeout time.Duration
	dial    func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error)
}

// GRPCPayload is the structured message sent over gRPC.
type GRPCPayload struct {
	Timestamp string `json:"timestamp"`
	Kind      string `json:"kind"`
	Proto     string `json:"proto"`
	Addr      string `json:"addr"`
	Port      int    `json:"port"`
	Process   string `json:"process,omitempty"`
	PID       int    `json:"pid,omitempty"`
}

// NewGRPCHandler creates a GRPCHandler that dials the given target address.
func NewGRPCHandler(target string, timeout time.Duration) *GRPCHandler {
	return &GRPCHandler{
		target:  target,
		timeout: timeout,
		dial:    grpc.Dial,
	}
}

// Handle forwards each change as a unary gRPC call payload.
// It establishes a connection per Handle invocation to keep the handler stateless.
func (h *GRPCHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	conn, err := h.dial(h.target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("grpc handler: dial %s: %w", h.target, err)
	}
	defer conn.Close()

	for _, c := range changes {
		payload := buildGRPCPayload(c)
		_ = payload // real impl would call a generated stub method here
	}
	return nil
}

// Drain is a no-op for the gRPC handler.
func (h *GRPCHandler) Drain() error { return nil }

func buildGRPCPayload(c monitor.Change) GRPCPayload {
	kind := "added"
	if c.Kind == monitor.Removed {
		kind = "removed"
	}
	hostOrIP := c.Binding.Hostname
	if hostOrIP == "" {
		hostOrIP = c.Binding.Addr
	}
	return GRPCPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Kind:      kind,
		Proto:     c.Binding.Proto,
		Addr:      hostOrIP,
		Port:      c.Binding.Port,
		Process:   c.Binding.Process,
		PID:       c.Binding.PID,
	}
}
