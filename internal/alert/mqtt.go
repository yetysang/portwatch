package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
)

// MQTTHandler publishes port change alerts to an MQTT broker.
type MQTTHandler struct {
	cfg    config.MQTTConfig
	client mqtt.Client
}

// NewMQTTHandler creates a new MQTTHandler and connects to the broker.
func NewMQTTHandler(cfg config.MQTTConfig) (*MQTTHandler, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(cfg.Broker).
		SetClientID(cfg.ClientID).
		SetConnectTimeout(cfg.Timeout).
		SetAutoReconnect(true)

	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
		opts.SetPassword(cfg.Password)
	}

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if !token.WaitTimeout(cfg.Timeout) {
		return nil, fmt.Errorf("mqtt: connection timed out to %s", cfg.Broker)
	}
	if err := token.Error(); err != nil {
		return nil, fmt.Errorf("mqtt: connect error: %w", err)
	}

	return &MQTTHandler{cfg: cfg, client: client}, nil
}

// Handle publishes each change as a JSON message to the configured topic.
func (h *MQTTHandler) Handle(_ context.Context, changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	for _, ch := range changes {
		payload, err := formatMQTTMsg(ch)
		if err != nil {
			return fmt.Errorf("mqtt: marshal error: %w", err)
		}
		token := h.client.Publish(h.cfg.Topic, h.cfg.QoS, h.cfg.Retain, payload)
		if !token.WaitTimeout(h.cfg.Timeout) {
			return fmt.Errorf("mqtt: publish timed out")
		}
		if err := token.Error(); err != nil {
			return fmt.Errorf("mqtt: publish error: %w", err)
		}
	}
	return nil
}

// Drain disconnects from the MQTT broker.
func (h *MQTTHandler) Drain() {
	h.client.Disconnect(250)
}

type mqttMessage struct {
	Action   string `json:"action"`
	Proto    string `json:"proto"`
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Process  string `json:"process,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	At       string `json:"at"`
}

func formatMQTTMsg(ch monitor.Change) ([]byte, error) {
	action := "added"
	if ch.Kind == monitor.Removed {
		action = "removed"
	}
	msg := mqttMessage{
		Action:   action,
		Proto:    ch.Binding.Proto,
		Addr:     ch.Binding.Addr,
		Port:     ch.Binding.Port,
		Process:  ch.Binding.Process,
		Hostname: ch.Binding.Hostname,
		At:       time.Now().UTC().Format(time.RFC3339),
	}
	return json.Marshal(msg)
}
