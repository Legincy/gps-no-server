package subscriptions

import (
	"context"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"gps-no-server/internal/logger"
	"gps-no-server/internal/models"
	"gps-no-server/internal/services"
	"time"
)

type RangingList []RangingItem

type RangingItem struct {
	SourceAddress      string       `json:"source_address"`
	DestinationAddress string       `json:"destination_address"`
	Distance           DistanceData `json:"distance"`
}

type DistanceData struct {
	RawDistance    float64 `json:"raw_distance"`
	ScaledDistance float64 `json:"scaled_distance"`
}

type RangingSubscription struct {
	log            zerolog.Logger
	rangingService *services.RangingService
	stationService *services.StationService
}

func NewRangingSubscription(rangingService *services.RangingService, stationService *services.StationService) *RangingSubscription {
	return &RangingSubscription{
		log:            logger.GetLogger("ranging-subscription"),
		rangingService: rangingService,
		stationService: stationService,
	}
}

func (c *RangingSubscription) GetTopics() []string {
	return []string{
		"gpsno/simulation/devices/+/uwb/ranging",
	}
}

func (c *RangingSubscription) HandleMessage(message mqtt.Message) {
	topic := message.Topic()
	payload := string(message.Payload())

	var rangingList RangingList
	if err := json.Unmarshal([]byte(payload), &rangingList); err != nil {
		c.log.Error().Err(err).Str("topic", topic).Msg("Failed to unmarshal message")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rangingModels []*models.Ranging
	for _, rangingData := range rangingList {
		source := &models.Station{MacAddress: rangingData.SourceAddress}
		destination := &models.Station{MacAddress: rangingData.DestinationAddress}

		sourceStation, _ := c.stationService.GetByMac(ctx, source.MacAddress)
		destinationStation, _ := c.stationService.GetByMac(ctx, destination.MacAddress)

		rangingModel := &models.Ranging{
			Source:      sourceStation,
			Destination: destinationStation,
			RawDistance: rangingData.Distance.RawDistance,
		}

		rangingModels = append(rangingModels, rangingModel)
	}

	if err := c.rangingService.SaveAll(ctx, rangingModels); err != nil {
		c.log.Error().Err(err).Msg("Failed to save rangingService data")
		return
	}
}
