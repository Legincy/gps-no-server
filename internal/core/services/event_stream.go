package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/core/models/dtos"
	"sync"
	"time"
)

type client struct {
	id      string
	channel chan []byte
	done    chan struct{}
}

type EventStreamService struct {
	clients        map[string]map[string]*client
	eventLock      sync.RWMutex
	log            zerolog.Logger
	maxClients     int
	clientTimeout  time.Duration
	messageTimeout time.Duration
}

func NewEventStreamService() *EventStreamService {
	return &EventStreamService{
		clients:        make(map[string]map[string]*client),
		log:            logger.GetLogger("event-stream-service"),
		maxClients:     1000,
		clientTimeout:  30 * time.Minute,
		messageTimeout: 5 * time.Second,
	}
}

func (s *EventStreamService) Subscribe(ctx context.Context, eventType string, filterId ...uint) (<-chan []byte, error) {
	subscriptionKey := eventType
	if len(filterId) > 0 && filterId[0] > 0 {
		subscriptionKey = fmt.Sprintf("%s:%d", eventType, filterId[0])
		s.log.Info().Str("key", subscriptionKey).Msg("Creating subscription for specific ID")
	}

	clientId := fmt.Sprintf("%s_%d", subscriptionKey, time.Now().UnixNano())

	s.eventLock.Lock()
	if s.getTotalClientCount() >= s.maxClients {
		s.eventLock.Unlock()
		return nil, fmt.Errorf("maximum number of clients reached")
	}

	messageChannel := make(chan []byte, 100)
	doneChannel := make(chan struct{})

	newClient := &client{
		id:      clientId,
		channel: messageChannel,
		done:    doneChannel,
	}

	if _, exists := s.clients[subscriptionKey]; !exists {
		s.clients[subscriptionKey] = make(map[string]*client)
	}
	s.clients[subscriptionKey][clientId] = newClient
	s.eventLock.Unlock()

	go func() {
		defer func() {
			s.removeClient(subscriptionKey, clientId)
			close(doneChannel)
			close(messageChannel)
		}()

		// Timeout-Timer starten
		timer := time.NewTimer(s.clientTimeout)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			s.log.Debug().Str("client", clientId).Msg("Client context cancelled")
		case <-timer.C:
			s.log.Warn().Str("client", clientId).Msg("Client timeout reached")
		case <-doneChannel:
			s.log.Debug().Str("client", clientId).Msg("Client done signal received")
		}
	}()

	return messageChannel, nil
}

func (s *EventStreamService) removeClient(subscriptionKey, clientId string) {
	s.eventLock.Lock()
	defer s.eventLock.Unlock()

	if clients, exists := s.clients[subscriptionKey]; exists {
		if client, exists := clients[clientId]; exists {
			select {
			case client.done <- struct{}{}:
			default:
			}
			delete(clients, clientId)
		}

		if len(clients) == 0 {
			delete(s.clients, subscriptionKey)
		}
	}
}

func (s *EventStreamService) getTotalClientCount() int {
	total := 0
	for _, clients := range s.clients {
		total += len(clients)
	}
	return total
}

func (s *EventStreamService) Publish(eventType string, data interface{}) error {
	var rangingID uint = 0

	if ranging, ok := data.(*dtos.RangingDto); ok && ranging.ID > 0 {
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
		return fmt.Errorf("error marshalling event data: %s", err)
	}
	message := []byte(fmt.Sprintf("event: %s\ndata: %s\n\n", eventType, jsonData))

	s.eventLock.RLock()
	defer s.eventLock.RUnlock()

	// Publish an alle Clients des Event-Typs
	if clients, exists := s.clients[eventType]; exists {
		s.publishToClients(eventType, clients, message)
	}

	// Publish an spezifische Ranging ID
	if rangingID > 0 {
		specificKey := fmt.Sprintf("%s:%d", eventType, rangingID)
		s.log.Debug().Str("key", specificKey).Msg("Publishing to specific ranging ID")

		if clients, exists := s.clients[specificKey]; exists {
			s.publishToClients(specificKey, clients, message)
		}
	}

	return nil
}

func (s *EventStreamService) publishToClients(key string, clients map[string]*client, message []byte) {
	var deadClients []string

	for clientId, client := range clients {
		select {
		case client.channel <- message:
		case <-time.After(s.messageTimeout):
			s.log.Warn().Str("client", clientId).Str("event_type", key).Msg("Client message timeout, marking for removal")
			deadClients = append(deadClients, clientId)
		case <-client.done:
			deadClients = append(deadClients, clientId)
		default:
			s.log.Warn().Str("client", clientId).Str("event_type", key).Msg("Client channel full, marking for removal")
			deadClients = append(deadClients, clientId)
		}
	}

	if len(deadClients) > 0 {
		go func() {
			s.eventLock.Lock()
			defer s.eventLock.Unlock()

			if clientsMap, stillExists := s.clients[key]; stillExists {
				for _, clientId := range deadClients {
					if client, exists := clientsMap[clientId]; exists {
						select {
						case client.done <- struct{}{}:
						default:
						}
						delete(clientsMap, clientId)
					}
				}

				if len(clientsMap) == 0 {
					delete(s.clients, key)
				}
			}
		}()
	}
}

func (s *EventStreamService) HandleSSERequest(c *gin.Context, eventType string, filterId ...uint) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Cache-Control")
	c.Writer.Flush()

	ctx, cancel := context.WithTimeout(c.Request.Context(), s.clientTimeout)
	defer cancel()

	eventChan, err := s.Subscribe(ctx, eventType, filterId...)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to subscribe to events")
		c.JSON(500, gin.H{"error": "Failed to subscribe to events"})
		return
	}

	clientGone := c.Writer.CloseNotify()
	go func() {
		select {
		case <-clientGone:
			s.log.Debug().Msg("Client disconnected")
			cancel()
		case <-ctx.Done():
			s.log.Debug().Msg("Context cancelled")
		}
	}()

	idInfo := ""
	if len(filterId) > 0 && filterId[0] > 0 {
		idInfo = fmt.Sprintf(`, "id": %d`, filterId[0])
	}

	connectMsg := fmt.Sprintf("event: connected\ndata: {\"status\":\"connected\"%s}\n\n", idInfo)
	if _, err := c.Writer.Write([]byte(connectMsg)); err != nil {
		s.log.Error().Err(err).Msg("Failed to write connection message")
		return
	}
	c.Writer.Flush()

	heartbeat := time.NewTicker(30 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case <-ctx.Done():
			s.log.Debug().Msg("SSE context done")
			return

		case <-heartbeat.C:
			heartbeatMsg := "event: heartbeat\ndata: {\"type\":\"heartbeat\"}\n\n"
			if _, err := c.Writer.Write([]byte(heartbeatMsg)); err != nil {
				s.log.Error().Err(err).Msg("Failed to write heartbeat")
				return
			}
			c.Writer.Flush()

		case message, ok := <-eventChan:
			if !ok {
				s.log.Debug().Msg("Event channel closed")
				return
			}

			if _, err := c.Writer.Write(message); err != nil {
				s.log.Error().Err(err).Msg("Failed to write SSE message")
				return
			}
			c.Writer.Flush()
		}
	}
}

func (s *EventStreamService) GetStats() map[string]interface{} {
	s.eventLock.RLock()
	defer s.eventLock.RUnlock()

	stats := make(map[string]interface{})
	stats["total_subscriptions"] = len(s.clients)
	stats["total_clients"] = s.getTotalClientCount()

	subscriptions := make(map[string]int)
	for key, clients := range s.clients {
		subscriptions[key] = len(clients)
	}
	stats["subscriptions"] = subscriptions

	return stats
}

func (s *EventStreamService) Shutdown(ctx context.Context) error {
	s.log.Info().Msg("Starting EventStreamService shutdown")

	s.eventLock.Lock()
	defer s.eventLock.Unlock()

	totalClients := s.getTotalClientCount()
	if totalClients > 0 {
		s.log.Info().Msgf("Shutting down %d active clients", totalClients)

		for eventType, clients := range s.clients {
			for clientId, client := range clients {
				s.log.Debug().
					Str("event_type", eventType).
					Str("client_id", clientId).
					Msg("Closing client connection")

				select {
				case client.done <- struct{}{}:
				default:
					// Channel ist bereits geschlossen oder voll
				}

				// Close channel, falls noch offen
				select {
				case <-client.channel:
				default:
					close(client.channel)
				}
			}
		}

		s.clients = make(map[string]map[string]*client)
		s.log.Info().Msg("All clients disconnected")
	}

	s.log.Info().Msg("EventStreamService shutdown completed")
	return nil
}

func (s *EventStreamService) IsHealthy() bool {
	s.eventLock.RLock()
	defer s.eventLock.RUnlock()

	totalClients := s.getTotalClientCount()
	return totalClients < s.maxClients
}
