package events

import (
	"sync"
	"time"
)

type StationEventType string

const (
	StationAddedToCluster     StationEventType = "station_added_to_cluster"
	StationRemovedFromCluster StationEventType = "station_removed_from_cluster"
)

type StationEvent struct {
	Type      StationEventType
	ClusterId uint
	StationId uint
	Timestamp time.Time
}

type StationEventBus struct {
	handlers map[StationEventType][]func(StationEvent)
	mu       sync.Mutex
}

func NewStationEventBus() *StationEventBus {
	return &StationEventBus{
		handlers: make(map[StationEventType][]func(StationEvent)),
	}
}

func (bus *StationEventBus) Subscribe(eventType StationEventType, handler func(StationEvent)) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if _, exists := bus.handlers[eventType]; !exists {
		bus.handlers[eventType] = []func(StationEvent){}
	}
	bus.handlers[eventType] = append(bus.handlers[eventType], handler)
}

func (bus *StationEventBus) Publish(event *StationEvent) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	eventX := *event

	if handlers, exists := bus.handlers[event.Type]; exists {
		for _, handler := range handlers {
			go handler(eventX)
		}
	}
}
