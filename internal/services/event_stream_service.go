package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gps-no-server/internal/controllers/dto"
	"gps-no-server/internal/logger"
	"sync"
)

type eventData struct {
	Type string
	Data interface{}
}

type EventStreamService struct {
	clients   map[string]map[chan []byte]bool
	eventLock sync.RWMutex
	log       zerolog.Logger
}

func NewEventStreamService() *EventStreamService {
	return &EventStreamService{
		clients: make(map[string]map[chan []byte]bool),
		log:     logger.GetLogger("event-stream-service"),
	}
}

func (s *EventStreamService) Subscribe(ctx context.Context, eventType string, filterId ...uint) (<-chan []byte, error) {
	s.eventLock.Lock()
	defer s.eventLock.Unlock()

	subscriptionKey := eventType
	if len(filterId) > 0 && filterId[0] > 0 {
		subscriptionKey = fmt.Sprintf("%s:%d", eventType, filterId[0])
		s.log.Info().Str("key", subscriptionKey).Msg("Creating subscription for specific ID")
	}

	messageChannel := make(chan []byte, 1024)

	if _, exists := s.clients[subscriptionKey]; !exists {
		s.clients[subscriptionKey] = make(map[chan []byte]bool)
	}
	s.clients[subscriptionKey][messageChannel] = true

	go func() {
		<-ctx.Done()
		s.eventLock.Lock()
		defer s.eventLock.Unlock()

		if channels, exists := s.clients[subscriptionKey]; exists {
			delete(channels, messageChannel)
			if len(channels) == 0 {
				delete(s.clients, subscriptionKey)
			}
		}

		close(messageChannel)
	}()

	return messageChannel, nil
}

func (s *EventStreamService) Unsubscribe(stationId string) error {
	return nil
}

func (s *EventStreamService) Publish(eventType string, data interface{}) error {
	var rangingID uint = 0

	if ranging, ok := data.(*dto.RangingDto); ok && ranging.ID > 0 {
		rangingID = ranging.ID
	} else if rangingMap, ok := data.(map[string]interface{}); ok {
		if id, exists := rangingMap["id"]; exists {
			if idFloat, ok := id.(float64); ok {
				rangingID = uint(idFloat)
			} else if idUint, ok := id.(uint); ok {
				rangingID = idUint
			}
		}
	} else if jsonStr, ok := data.(string); ok {
		var jsonMap map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &jsonMap); err == nil {
			if id, exists := jsonMap["id"]; exists {
				if idFloat, ok := id.(float64); ok {
					rangingID = uint(idFloat)
				}
			}
		}
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Error marshalling event data: %s", err)
	}
	message := []byte(fmt.Sprintf("event : %s\ndata: %s\n\n", eventType, jsonData))

	s.eventLock.RLock()
	defer s.eventLock.RUnlock()

	// First, publish to the general channel
	if channels, exists := s.clients[eventType]; exists {
		s.publishToChannels(eventType, channels, message)
	}

	// If we have a valid ID, also publish to the specific ID channel
	if rangingID > 0 {
		specificKey := fmt.Sprintf("%s:%d", eventType, rangingID)
		s.log.Debug().Str("key", specificKey).Msg("Publishing to specific ranging ID")

		if channels, exists := s.clients[specificKey]; exists {
			s.publishToChannels(specificKey, channels, message)
		}
	}

	return nil
}

func (s *EventStreamService) publishToChannels(key string, channels map[chan []byte]bool, message []byte) {
	var slowClients []chan []byte

	for ch := range channels {
		select {
		case ch <- message:
		default:
			s.log.Warn().Str("event_type", key).Msg("Channel is full, skipping message")
			slowClients = append(slowClients, ch)
		}
	}

	if len(slowClients) > 0 {
		s.eventLock.Lock()
		if channelsMap, stillExists := s.clients[key]; stillExists {
			for _, ch := range slowClients {
				delete(channelsMap, ch)
				close(ch)
			}

			if len(channelsMap) == 0 {
				delete(s.clients, key)
			}
		}
		s.eventLock.Unlock()
	}
}

func (s *EventStreamService) HandleSSERequest(c *gin.Context, eventType string, filterId ...uint) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Flush()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	eventChan, err := s.Subscribe(ctx, eventType, filterId...)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	clientGone := c.Writer.CloseNotify()
	go func() {
		<-clientGone
		cancel()
	}()

	idInfo := ""
	if len(filterId) > 0 && filterId[0] > 0 {
		idInfo = fmt.Sprintf(`, "id": %d`, filterId[0])
	}
	fmt.Fprintf(c.Writer, "event: connected\ndata: {\"status\":\"connected\"%s}\n\n", idInfo)
	c.Writer.Flush()

	for {
		select {
		case <-ctx.Done():
			return
		case message, ok := <-eventChan:
			if !ok {
				return
			}
			_, err := c.Writer.Write(message)
			if err != nil {
				s.log.Error().Err(err).Msg("Failed to write SSE message")
				return
			}
			c.Writer.Flush()
		}
	}
}

func (s *EventStreamService) mergeChannels(ctx context.Context, channels ...<-chan []byte) <-chan []byte {
	out := make(chan []byte)

	var wg sync.WaitGroup
	wg.Add(len(channels))

	for _, c := range channels {
		go func(c <-chan []byte) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-c:
					if !ok {
						return
					}
					select {
					case out <- msg:
					case <-ctx.Done():
						return
					}
				}
			}
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
