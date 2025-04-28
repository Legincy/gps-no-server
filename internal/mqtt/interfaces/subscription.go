package interfaces

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Subscription interface {
	GetTopics() []string
	HandleMessage(message mqtt.Message)
}
