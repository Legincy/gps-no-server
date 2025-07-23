package interfaces

import "gorm.io/gorm"

type DatabaseConnection interface {
	GetDB() *gorm.DB
	Close() error
	Migrate() error
}

type MQTTConnection interface {
	Connect() error
	Disconnect() error
	Subscribe(topic string, handler func([]byte)) error
	Publish(topic string, payload []byte) error
	IsConnected() bool
}

type InfluxConnection interface {
	Connect() error
	Disconnect() error
	WritePoint(measurement string, tags map[string]string, fields map[string]interface{}) error
	IsConnected() bool
}
