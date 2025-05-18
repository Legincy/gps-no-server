package interfaces

import (
	"context"
	"github.com/gin-gonic/gin"
)

type EventStreamSubscriber interface {
	Subscribe(ctx context.Context, stationId uint) (<-chan interface{}, error)
	Unsubscribe(stationId uint)
}

type EventStreamPublisher interface {
	Publish(eventType string, data interface{}) error
}

type EventStreamService interface {
	EventStreamSubscriber
	EventStreamPublisher
	HandleSSERequest(c *gin.Context, eventType string, filterId ...uint)
}
