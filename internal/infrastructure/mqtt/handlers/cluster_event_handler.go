package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/events"
	"gps-no-server/internal/infrastructure/mqtt"
)

type ClusterEventHandler struct {
	mqttClient *mqtt.Client
	log        zerolog.Logger
}

func NewClusterEventHandler(mqttClient *mqtt.Client) *ClusterEventHandler {
	return &ClusterEventHandler{
		mqttClient: mqttClient,
		log:        logger.GetLogger("cluster-event-handler"),
	}
}

func (c *ClusterEventHandler) HandleEvent(event *events.StationEvent) {
	topic := c.buildTopic(event)
	payload, err := json.Marshal(event)

	if err != nil {
		c.log.Error().Err(err).Msg("Failed to marshal cluster event")
		return
	}

	/*
		if err := c.mqttClient.Publish(topic, 1, false, payload); err != nil {
			c.log.Error().Str("topic", topic).Msg("Failed to publish event")
			return
		}
	*/

	c.log.Info().
		Str("type", string(event.Type)).
		Str("topic", topic).
		Str("payload", string(payload)).
		Msg("Published cluster event to MQTT")
}

func (c *ClusterEventHandler) buildTopic(event *events.StationEvent) string {
	switch event.Type {
	case events.StationAddedToCluster, events.StationRemovedFromCluster:
		return fmt.Sprintf("gpsno/clusters/%d/stations/%d", event.ClusterId, event.StationId)
	case events.ClusterUpdated:
		return fmt.Sprintf("gpsno/clusters/%d", event.ClusterId)
	default:
		return fmt.Sprintf("gpsno/clusters/events")
	}
}
