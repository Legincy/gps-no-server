package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gps-no-server/internal/common/config"
	"gps-no-server/internal/common/logger"
	"math/rand"
	"time"
)

type Client struct {
	client   mqtt.Client
	config   *config.MqttConfig
	Registry *Registry
	log      zerolog.Logger
}

func Create(cfg *config.MqttConfig, registry *Registry) (*Client, error) {
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
		client:   mqtt.NewClient(opts),
		config:   cfg,
		log:      logger.GetLogger("mqtt"),
		Registry: registry,
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

func (c *Client) SubscribeRegistry() error {
	topics := c.Registry.GetAllTopics()

	for _, topic := range topics {
		if err := c.Subscribe(topic, c.messageHandler); err != nil {
			return fmt.Errorf("Failed to subscribe to topic %s: %w", topic, err)
		}
		c.log.Info().Msgf("Subscribed to topic %s", topic)
	}

	return nil
}

func (c *Client) messageHandler(client mqtt.Client, message mqtt.Message) {
	topic := message.Topic()
	handlers := c.Registry.GetAllSubscriptions(topic)

	if len(handlers) > 0 {
		for _, handler := range handlers {
			handler.HandleMessage(message)
		}
	} else {
		log.Warn().Msgf("No handler found for topic %s", topic)
	}
}

func (c *Client) Publish(topic string, qos int, retained bool, payload []byte) interface{} {
	if token := c.client.Publish(topic, byte(qos), retained, payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
