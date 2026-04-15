package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// RedisPublisher is the interface for publishing messages to Redis streams.
type RedisPublisher interface {
	XAdd(ctx context.Context, stream, id string, values map[string]interface{}) error
	Close() error
}

// RedisHandler publishes port change events to a Redis stream.
type RedisHandler struct {
	client RedisPublisher
	stream string
}

// NewRedisHandler creates a RedisHandler that publishes to the given stream.
func NewRedisHandler(client RedisPublisher, stream string) *RedisHandler {
	return &RedisHandler{client: client, stream: stream}
}

func (h *RedisHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, c := range changes {
		msg := formatRedisMsg(c)
		if err := h.client.XAdd(ctx, h.stream, "*", msg); err != nil {
			return fmt.Errorf("redis xadd: %w", err)
		}
	}
	return nil
}

func (h *RedisHandler) Drain() error {
	return nil
}

func formatRedisMsg(c monitor.Change) map[string]interface{} {
	payload, _ := json.Marshal(map[string]interface{}{
		"kind":     string(c.Kind),
		"proto":    c.Binding.Proto,
		"port":     c.Binding.Port,
		"addr":     c.Binding.Addr,
		"pid":      c.Binding.PID,
		"process":  c.Binding.Process,
		"hostname": c.Binding.Hostname,
	})
	return map[string]interface{}{
		"event": string(payload),
		"ts":    time.Now().UTC().Format(time.RFC3339),
	}
}
