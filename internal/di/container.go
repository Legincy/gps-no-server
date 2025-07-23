package di

import (
	"fmt"
	"gps-no-server/internal/common/config"
	"gps-no-server/internal/core/repositories"
	"gps-no-server/internal/core/services"
	"gps-no-server/internal/di/interfaces"
	"gps-no-server/internal/infrastructure/database"
	"gps-no-server/internal/infrastructure/database/influx"
	"gps-no-server/internal/infrastructure/mqtt"
	"gps-no-server/internal/infrastructure/mqtt/handlers"
)

type Container struct {
	config *config.Config

	database interfaces.DatabaseConnection
	mqtt     interfaces.MQTTConnection
	influx   interfaces.InfluxConnection

	stationRepo interfaces.StationRepository
	clusterRepo interfaces.ClusterRepository
	rangingRepo interfaces.RangingRepository

	stationService interfaces.StationService
	clusterService interfaces.ClusterService
	rangingService interfaces.RangingService

	stationHandler interfaces.StationHandler
	clusterHandler interfaces.ClusterHandler
	rangingHandler interfaces.RangingHandler

	mqttStationHandler interfaces.MQTTStationHandler
	mqttRangingHandler interfaces.MQTTRangingHandler
}

func NewContainer(cfg *config.Config) (*Container, error) {
	c := &Container{
		config: cfg,
	}

	if err := c.initInfrastructure(); err != nil {
		return nil, fmt.Errorf("failed to initialize infrastructure: %w", err)
	}

	if err := c.initRepositories(); err != nil {
		return nil, fmt.Errorf("failed to initialize repositories: %w", err)
	}

	if err := c.initServices(); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}

	if err := c.initHandlers(); err != nil {
		return nil, fmt.Errorf("failed to initialize handlers: %w", err)
	}

	/*
		if err := c.initMQTTSubscriptions(); err != nil {
			return nil, fmt.Errorf("failed to initialize MQTT subscriptions: %w", err)
		}
	*/

	return c, nil
}

func (c *Container) initInfrastructure() error {
	db, err := database.NewPostgresConnection(&c.config.Database)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	c.database = db

	if err := c.database.Migrate(); err != nil {
		return fmt.Errorf("failed to run database migrations: %w", err)
	}

	mqttClient, err := mqtt.NewMQTTConnection(&c.config.Mqtt)
	if err != nil {
		return fmt.Errorf("failed to initialize MQTT: %w", err)
	}
	c.mqtt = mqttClient

	influxClient, err := influx.NewInfluxConnection(&c.config.Influx)
	if err != nil {
		return fmt.Errorf("failed to initialize InfluxDB: %w", err)
	}
	c.influx = influxClient

	return nil
}

func (c *Container) initRepositories() error {
	c.stationRepo = repositories.NewStationRepository(c.database.GetDB())
	c.clusterRepo = repositories.NewClusterRepository(c.database.GetDB())
	c.rangingRepo = repositories.NewRangingRepository(c.database.GetDB())

	return nil
}

func (c *Container) initServices() error {
	c.stationService = services.NewStationService(c.stationRepo)
	c.clusterService = services.NewClusterService(c.clusterRepo)
	c.rangingService = services.NewRangingService(c.rangingRepo, c.influx)

	return nil
}

func (c *Container) initHandlers() error {
	c.stationHandler = handlers.NewStationHandler(c.stationService)
	//c.clusterHandler = handlers.NewClusterHandler(c.clusterService)
	//c.rangingHandler = handlers.NewRangingHandler(c.rangingService)

	//c.mqttStationHandler = mqtthandlers.NewStationHandler(c.stationService)
	//c.mqttRangingHandler = mqtthandlers.NewRangingHandler(c.rangingService)

	return nil
}

func (c *Container) initMQTTSubscriptions() error {
	if err := c.mqtt.Subscribe("gpsno/simulation/devices/+/device/raw",
		func(payload []byte) {
			if err := c.mqttStationHandler.HandleStationData(payload); err != nil {
			}
		}); err != nil {
		return fmt.Errorf("failed to subscribe to station topic: %w", err)
	}

	if err := c.mqtt.Subscribe("gpsno/simulation/devices/+/uwb/ranging",
		func(payload []byte) {
			if err := c.mqttRangingHandler.HandleRangingData(payload); err != nil {
			}
		}); err != nil {
		return fmt.Errorf("failed to subscribe to ranging topic: %w", err)
	}

	return nil
}

func (c *Container) GetStationHandler() interfaces.StationHandler {
	fmt.Printf("%+v\n", c)
	return c.stationHandler
}

func (c *Container) GetClusterHandler() interfaces.ClusterHandler {
	return c.clusterHandler
}

func (c *Container) GetRangingHandler() interfaces.RangingHandler {
	return c.rangingHandler
}

func (c *Container) GetMQTTConnection() interfaces.MQTTConnection {
	return c.mqtt
}

func (c *Container) GetInfluxConnection() interfaces.InfluxConnection {
	return c.influx
}

func (c *Container) Cleanup() error {
	var errors []error

	if c.mqtt != nil && c.mqtt.IsConnected() {
		if err := c.mqtt.Disconnect(); err != nil {
			errors = append(errors, fmt.Errorf("failed to disconnect MQTT: %w", err))
		}
	}

	if c.influx != nil && c.influx.IsConnected() {
		if err := c.influx.Disconnect(); err != nil {
			errors = append(errors, fmt.Errorf("failed to disconnect InfluxDB: %w", err))
		}
	}

	if c.database != nil {
		if err := c.database.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close database: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("cleanup errors: %v", errors)
	}

	return nil
}
