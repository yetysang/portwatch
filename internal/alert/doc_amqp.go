// Package alert provides alert handlers for portwatch.
//
// # AMQP Handler
//
// The AMQPHandler publishes port change events to an AMQP 0-9-1 compatible
// broker such as RabbitMQ. Each change is serialised as a JSON object and
// delivered to a configurable exchange with a routing key.
//
// Example message payload:
//
//	{
//	  "timestamp": "2024-01-15T10:30:00Z",
//	  "action":    "added",
//	  "proto":     "tcp",
//	  "addr":      "0.0.0.0",
//	  "port":      8080,
//	  "process":   "nginx",
//	  "pid":       1234
//	}
//
// Usage:
//
//	pub := amqpbroker.Dial(cfg.AMQP.URL)
//	h := alert.NewAMQPHandler(pub, cfg.AMQP.Exchange, cfg.AMQP.RoutingKey)
//
// The AMQPPublisher interface allows the concrete broker connection to be
// swapped for a stub in tests without importing a specific AMQP library.
package alert
