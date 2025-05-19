package di

import (
	"github.com/rs/zerolog/log"
	"gps-no-server/internal/common/config"
	"gps-no-server/internal/core/controllers"
	"gps-no-server/internal/core/repositories"
	"gps-no-server/internal/core/services"
	"gps-no-server/internal/events"
	"gps-no-server/internal/infrastructure/database"
	"gps-no-server/internal/infrastructure/mqtt"
	"gps-no-server/internal/infrastructure/mqtt/handlers"
	"gps-no-server/internal/infrastructure/mqtt/subscriptions"
)

type Container struct {
	Config             *config.Config
	Database           *database.GormDB
	MqttClient         *mqtt.Client
	EventStreamService *services.EventStreamService

	StationEventBus     *events.StationEventBus
	ClusterEventHandler *handlers.ClusterEventHandler

	StationRepository *repositories.StationRepository
	ClusterRepository *repositories.ClusterRepository
	RangingRepository *repositories.RangingRepository

	StationService *services.StationService
	ClusterService *services.ClusterService
	RangingService *services.RangingService

	StationController *controllers.StationController
	ClusterController *controllers.ClusterController
	RangingController *controllers.RangingController
	EventController   *controllers.EventStreamController
}

func NewContainer(cfg *config.Config) (*Container, error) {
	container := &Container{
		Config: cfg,
	}

	if err := container.initDatabase(); err != nil {
		return nil, err
	}

	container.initRepositories()
	container.initServices()
	container.initControllers()
	container.initMqtt()
	container.initEvents()

	return container, nil
}

func (c *Container) initDatabase() error {
	db, err := database.NewGormDB(&c.Config.Database)
	if err != nil {
		return err
	}
	c.Database = db

	return nil
}

func (c *Container) initRepositories() {
	c.StationRepository = repositories.NewStationRepository(c.Database.DB)
	c.ClusterRepository = repositories.NewClusterRepository(c.Database.DB)
	c.RangingRepository = repositories.NewRangingRepository(c.Database.DB)
}

func (c *Container) initServices() {
	c.StationService = services.NewStationService(c.StationRepository, c.StationEventBus)
	c.ClusterService = services.NewClusterService(c.ClusterRepository)
	c.RangingService = services.NewRangingService(c.RangingRepository, c.StationService, c.EventStreamService)
	c.EventStreamService = services.NewEventStreamService()
}

func (c *Container) initControllers() {
	c.StationController = controllers.NewStationController(c.StationService)
	c.ClusterController = controllers.NewClusterController(c.ClusterService)
	c.RangingController = controllers.NewRangingController(c.RangingService)
	c.EventController = controllers.NewEventStreamController(c.EventStreamService)
}

func (c *Container) initEvents() {
	c.StationEventBus = events.NewStationEventBus()

	clusterEventHandler := handlers.NewClusterEventHandler(c.MqttClient)

	c.StationEventBus.Subscribe(events.StationAddedToCluster, func(event events.StationEvent) {
		clusterEventHandler.HandleEvent(&event)
	})
}

func (c *Container) initMqtt() {
	mqttRegistry := mqtt.NewSubscriptionRegistry()

	stationHandler := subscriptions.NewStationSubscription(c.StationService)
	rangingHandler := subscriptions.NewRangingSubscription(c.RangingService, c.StationService)

	mqttRegistry.Register(stationHandler)
	mqttRegistry.Register(rangingHandler)

	mqttClient, _ := mqtt.Create(&c.Config.Mqtt, mqttRegistry)
	c.MqttClient = mqttClient
}

func (c *Container) Cleanup() {
	if err := c.Database.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close database connection")
	}

	if err := c.MqttClient.Disconnect(); err != nil {
		log.Error().Err(err).Msg("Failed to disconnect MQTT client")
	}
}
