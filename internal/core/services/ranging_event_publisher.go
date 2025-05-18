package services

import (
	"context"
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/core/models/mappers"
)

const RangingEventType = "update"

type RangingEventPublisher struct {
	eventService *EventStreamService
}

func NewRangingEventPublisher(eventService *EventStreamService) *RangingEventPublisher {
	return &RangingEventPublisher{
		eventService: eventService,
	}
}

func (p *RangingEventPublisher) PublishRangingEvent(ctx context.Context, ranging *models.Ranging) error {
	rangingDto := mappers.FromRanging(ranging, nil)

	return p.eventService.Publish(RangingEventType, rangingDto)
}
