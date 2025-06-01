package di

import (
	"context"
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
	"sync"
	"time"
)

type Container struct {
	Config             *config.Config
	Database           *database.GormDB
	MqttClient         *mqtt.Client
	EventStreamService *services.EventStreamService

	StationEventBus     *events.StationEventBus
	ClusterEventHandler *handlers.ClusterEventHandler

	StationRepository       *repositories.StationRepository
	StationConfigRepository *repositories.StationConfigurationRepository
	ClusterRepository       *repositories.ClusterRepository
	RangingRepository       *repositories.RangingRepository

	StationService       *services.StationService
	StationConfigService *services.StationConfigurationService
	ClusterService       *services.ClusterService
	RangingService       *services.RangingService

	StationController       *controllers.StationController
	StationConfigController *controllers.StationConfigController
	ClusterController       *controllers.ClusterController
	RangingController       *controllers.RangingController

	cleanupCtx     context.Context
	cleanupCancel  context.CancelFunc
	cleanupWg      sync.WaitGroup
	isShuttingDown bool
	shutdownMu     sync.RWMutex
}

func NewContainer(cfg *config.Config) (*Container, error) {
	cleanupCtx, cleanupCancel := context.WithCancel(context.Background())

	container := &Container{
		Config:        cfg,
		cleanupCtx:    cleanupCtx,
		cleanupCancel: cleanupCancel,
	}

	container.initEventStream()

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

func (c *Container) initEventStream() {
	c.EventStreamService = services.NewEventStreamService()
	log.Info().Msg("EventStreamService initialized")
}

func (c *Container) initDatabase() error {
	maxRetries := 5
	var err error

	for i := 0; i < maxRetries; i++ {
		c.Database, err = database.NewGormDB(&c.Config.Database)
		if err == nil {
			log.Info().Msg("Database connected successfully")
			return nil
		}

		log.Warn().Err(err).Msgf("Database connection attempt %d/%d failed", i+1, maxRetries)
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}

	return err
}

func (c *Container) initRepositories() {
	c.StationRepository = repositories.NewStationRepository(c.Database.DB)
	c.StationConfigRepository = repositories.NewStationConfigRepository(c.Database.DB)
	c.ClusterRepository = repositories.NewClusterRepository(c.Database.DB)
	c.RangingRepository = repositories.NewRangingRepository(c.Database.DB)
	log.Info().Msg("Repositories initialized")
}

func (c *Container) initServices() {
	c.StationService = services.NewStationService(c.StationRepository)
	c.StationConfigService = services.NewStationConfigService(c.StationConfigRepository)
	c.ClusterService = services.NewClusterService(c.ClusterRepository)
	c.RangingService = services.NewRangingService(c.RangingRepository, c.StationService, c.EventStreamService)
	log.Info().Msg("Services initialized")

}

func (c *Container) initControllers() {
	c.StationController = controllers.NewStationController(c.StationService)
	c.StationConfigController = controllers.NewStationConfigController(c.StationConfigService)
	c.ClusterController = controllers.NewClusterController(c.ClusterService)
	c.RangingController = controllers.NewRangingController(c.RangingService, c.EventStreamService)
	log.Info().Msg("Controllers initialized")
}

func (c *Container) initEvents() {
	c.StationEventBus = events.NewStationEventBus()

	c.ClusterEventHandler = handlers.NewClusterEventHandler(c.MqttClient)

	c.StationEventBus.Subscribe(events.StationAddedToCluster, func(event events.StationEvent) {
		c.shutdownMu.RLock()
		defer c.shutdownMu.RUnlock()

		if !c.isShuttingDown {
			c.ClusterEventHandler.HandleEvent(&event)
		}
	})
	log.Info().Msg("Event handlers initialized")
}

func (c *Container) initMqtt() {
	mqttRegistry := mqtt.NewSubscriptionRegistry()

	stationHandler := subscriptions.NewStationSubscription(c.StationService)
	rangingHandler := subscriptions.NewRangingSubscription(c.RangingService, c.StationService)

	mqttRegistry.Register(stationHandler)
	mqttRegistry.Register(rangingHandler)

	mqttClient, err := mqtt.Create(&c.Config.Mqtt, mqttRegistry)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create MQTT client")
	} else {
		c.MqttClient = mqttClient
		log.Info().Msg("MQTT client initialized")
	}
}

func (c *Container) IsShuttingDown() bool {
	c.shutdownMu.RLock()
	defer c.shutdownMu.RUnlock()
	return c.isShuttingDown
}

func (c *Container) GetCleanupContext() context.Context {
	return c.cleanupCtx
}

func (c *Container) Cleanup() {
	log.Info().Msg("Starting cleanup process")

	c.shutdownMu.Lock()
	c.isShuttingDown = true
	c.shutdownMu.Unlock()

	c.cleanupCancel()

	cleanupTimeout := 30 * time.Second
	cleanupCtx, cancel := context.WithTimeout(context.Background(), cleanupTimeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		c.cleanupWg.Wait()
	}()

	c.cleanupWg.Add(1)
	go func() {
		defer c.cleanupWg.Done()
		c.cleanupEventStream()
	}()

	if c.MqttClient != nil {
		c.cleanupWg.Add(1)
		go func() {
			defer c.cleanupWg.Done()
			c.cleanupMqtt()
		}()
	}

	if c.Database != nil {
		c.cleanupWg.Add(1)
		go func() {
			defer c.cleanupWg.Done()
			c.cleanupDatabase()
		}()
	}

	select {
	case <-done:
		log.Info().Msg("All cleanup operations completed successfully")
	case <-cleanupCtx.Done():
		log.Warn().Msg("Cleanup timeout reached, forcing shutdown")
	}

	log.Info().Msg("Cleanup process completed")
}

func (c *Container) cleanupEventStream() {
	if c.EventStreamService != nil {
		log.Info().Msg("Cleaning up EventStreamService...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := c.EventStreamService.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("Error during EventStreamService shutdown")
		} else {
			log.Info().Msg("EventStreamService shutdown completed successfully")
		}
	}
}

func (c *Container) cleanupMqtt() {
	log.Info().Msg("Cleaning up MQTT client...")

	disconnectTimeout := 5 * time.Second
	disconnectCtx, cancel := context.WithTimeout(context.Background(), disconnectTimeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- c.MqttClient.Disconnect()
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Error().Err(err).Msg("Error during MQTT disconnect")
		} else {
			log.Info().Msg("MQTT client disconnected successfully")
		}
	case <-disconnectCtx.Done():
		log.Warn().Msg("MQTT disconnect timeout, forcing close")
	}
}

func (c *Container) cleanupDatabase() {
	log.Info().Msg("Cleaning up database connection...")

	dbCloseTimeout := 10 * time.Second
	dbCloseCtx, cancel := context.WithTimeout(context.Background(), dbCloseTimeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- c.Database.Close()
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Error().Err(err).Msg("Error closing database connection")
		} else {
			log.Info().Msg("Database connection closed successfully")
		}
	case <-dbCloseCtx.Done():
		log.Warn().Msg("Database close timeout reached")
	}
}
