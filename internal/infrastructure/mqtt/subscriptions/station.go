package subscriptions

import (
	"context"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/core/services"
	"time"
)

type StationRaw struct {
	UWB struct {
		DeviceType string `json:"device_type"`
		Cluster    struct {
			Name     *string  `json:"name"`
			Stations []string `json:"stations"`
		} `json:"cluster"`
	} `json:"uwb"`
	Device struct {
		MacAddress string  `json:"mac_address"`
		Name       string  `json:"name"`
		Randomizer float64 `json:"randomizer"`
		UpdatedAt  string  `json:"updated_at"`
		CreatedAt  string  `json:"created_at"`
		StartedAt  string  `json:"started_at"`
		Uptime     int64   `json:"uptime"`
		Position   struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"position"`
	} `json:"device"`
}

type StationSubscription struct {
	log            zerolog.Logger
	stationService *services.StationService
}

func NewStationSubscription(stationService *services.StationService) *StationSubscription {
	return &StationSubscription{
		log:            logger.GetLogger("station-subscription"),
		stationService: stationService,
	}
}

func (c *StationSubscription) GetTopics() []string {
	return []string{
		"gpsno/simulation/devices/+/device/raw",
	}
}

func (c *StationSubscription) HandleMessage(message mqtt.Message) {
	topic := message.Topic()
	payload := string(message.Payload())

	var stationRaw StationRaw
	if err := json.Unmarshal([]byte(payload), &stationRaw); err != nil {
		c.log.Error().Err(err).Str("topic", topic).Msg("Failed to unmarshal station raw data")
	}

	station := &models.Station{
		MacAddress: stationRaw.Device.MacAddress,
		Name:       stationRaw.Device.Name,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := c.stationService.UpdateOrCreate(ctx, station, nil); err != nil {
		c.log.Error().Err(err).Str("mac", station.MacAddress).Msg("Failed to save station")
		return
	}

	c.log.Debug().Str("mac", station.MacAddress).Msg("Station data saved successfully")

}
