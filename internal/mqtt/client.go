package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"gps-no-server/internal/config"
	"gps-no-server/internal/logger"
	"math/rand"
	"time"
)

var log zerolog.Logger

type Client struct {
	client mqtt.Client
	config *config.MqttConfig
}

func Create(cfg *config.MqttConfig) (*Client, error) {
	log = logger.GetLogger("mqtt")

	opts := mqtt.NewClientOptions()
	broker := fmt.Sprintf("tcp://%s:%d", cfg.Host, cfg.Port)
	randomInt := rand.Intn(16777216)
	clientId := fmt.Sprintf("%s-%06x", cfg.ClientId, randomInt)

	opts.AddBroker(broker)
	opts.SetClientID(clientId)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(1 * time.Second)
	opts.SetCleanSession(true)

	if cfg.Username != "" && cfg.Password != "" {
		opts.SetUsername(cfg.Username)
		opts.SetPassword(cfg.Password)
	}

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Info().Msgf("Successfully connected to MQTT broker: %s", broker)
	})

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Info().Msgf("Connection lost: %v", err)
	})

	mqttClient := &Client{
		client: mqtt.NewClient(opts),
		config: cfg,
	}

	return mqttClient, nil
}

func (c *Client) Subscribe(topic string, callback mqtt.MessageHandler) error {
	if token := c.client.Subscribe(topic, 0, callback); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (c *Client) DefaultCallback(client mqtt.Client, msg mqtt.Message) {
	log.Info().Msgf("Received message on topic %s: %s", msg.Topic(), string(msg.Payload()))
}

func (c *Client) Connect() error {
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *Client) Disconnect() {
	if !c.client.IsConnected() {
		return
	}

	c.client.Disconnect(250)
}
