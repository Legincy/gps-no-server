package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gps-no-server/internal/common/config"
	"math/rand"
)

type MQTTConnection struct {
	client mqtt.Client
	config *config.MqttConfig
}

func NewMQTTConnection(cfg *config.MqttConfig) (*MQTTConnection, error) {
	opts := mqtt.NewClientOptions()
	broker := fmt.Sprintf("tcp://%s:%d", cfg.Host, cfg.Port)
	randomInt := rand.Intn(16777216)
	clientId := fmt.Sprintf("%s-%06x", cfg.ClientId, randomInt)

	opts.AddBroker(broker)
	opts.SetClientID(clientId)
	opts.SetAutoReconnect(cfg.AutoReconnect)
	opts.SetMaxReconnectInterval(cfg.MaxReconnectInterval)
	opts.SetCleanSession(cfg.CleanSession)

	if cfg.Username != "" && cfg.Password != "" {
		opts.SetUsername(cfg.Username)
		opts.SetPassword(cfg.Password)
	}

	client := mqtt.NewClient(opts)

	return &MQTTConnection{
		client: client,
		config: cfg,
	}, nil
}

func (m *MQTTConnection) Connect() error {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MQTTConnection) Disconnect() error {
	if m.client.IsConnected() {
		m.client.Disconnect(250)
	}
	return nil
}

func (m *MQTTConnection) Subscribe(topic string, handler func([]byte)) error {
	callback := func(client mqtt.Client, msg mqtt.Message) {
		handler(msg.Payload())
	}

	if token := m.client.Subscribe(topic, 0, callback); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MQTTConnection) Publish(topic string, payload []byte) error {
	if token := m.client.Publish(topic, 0, false, payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *MQTTConnection) IsConnected() bool {
	return m.client.IsConnected()
}
