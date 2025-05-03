package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gps-no-server/internal/cache"
	"gps-no-server/internal/config"
	"gps-no-server/internal/controllers"
	"gps-no-server/internal/database"
	"gps-no-server/internal/logger"
	"gps-no-server/internal/mqtt"
	"gps-no-server/internal/mqtt/subscriptions"
	"gps-no-server/internal/repository"
	"gps-no-server/internal/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logLevel := cfg.Server.LogLevel
	logger.Init(logLevel)
	log := logger.GetLogger("main")

	gormDB, err := database.NewGormDB(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database connection")
	}
	defer func() {
		if err := gormDB.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing database connection")
		}
	}()

	cacheManager := cache.NewCacheManager()

	stationRepository := repository.NewStationRepository(gormDB.DB)
	rangingRepository := repository.NewRangingRepository(gormDB.DB)
	clusterRepository := repository.NewClusterRepository(gormDB.DB)

	stationService := services.NewStationService(stationRepository)
	rangingService := services.NewRangingService(rangingRepository, stationRepository, cacheManager)
	clusterService := services.NewClusterService(clusterRepository)

	mqttClient, err := initMqtt(&cfg.Mqtt, stationService, rangingService)
	if err != nil {
		log.Error().Err(err).Msg("Error while initializing MQTT client")
	}
	defer mqttClient.Disconnect()

	server, err := initServer(&cfg.Server, stationService, rangingService, clusterService)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while initializing server")
	}

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Info().Msg("Shutdown signal received")

	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout*time.Second)
		defer cancel()

		log.Info().Msg("Shutting down server...")
		if err := server.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Error while shutting down server")
		}
		log.Info().Msg("Successfully stopped HTTP server")
	}
}

func initServer(
	cfg *config.ServerConfig,
	stationService *services.StationService,
	rangingService *services.RangingService,
	clusterService *services.ClusterService,
) (*http.Server, error) {
	gin.SetMode(cfg.ReleaseMode)
	router := gin.New()

	if cfg.LogLevel == "debug" {
		router.Use(gin.Logger())
	}

	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	stationController := controllers.NewStationController(stationService)
	clusterController := controllers.NewClusterController(clusterService)
	rangingController := controllers.NewRangingController(rangingService)

	apiHandler := controllers.NewAPI(stationController, clusterController, rangingController)
	apiHandler.RegisterRoutes(router)

	// Server erstellen
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout * time.Second,
		WriteTimeout: cfg.WriteTimeout * time.Second,
		IdleTimeout:  120 * time.Second, // Zusätzlicher Idle-Timeout für Keep-Alive-Verbindungen
	}

	// Server starten
	errChan := make(chan error, 1)

	go func() {
		log.Info().Msgf("Starting server on %s:%d (%s)", cfg.Host, cfg.Port, cfg.ReleaseMode)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Server failed to start")
			errChan <- err
		}
	}()

	// Kurze Wartezeit, um sicherzustellen, dass der Server starten konnte
	select {
	case err := <-errChan:
		return nil, err
	case <-time.After(100 * time.Millisecond):
		log.Info().Msgf("Server started successfully on %s:%d (%s)", cfg.Host, cfg.Port, cfg.ReleaseMode)
	}

	return server, nil
}

func initMqtt(
	cfg *config.MqttConfig,
	stationService *services.StationService,
	rangingService *services.RangingService,
) (*mqtt.Client, error) {
	subscriptionRegistry := mqtt.NewSubscriptionRegistry()

	stationHandler := subscriptions.NewStationSubscription(stationService)
	rangingHandler := subscriptions.NewRangingSubscription(rangingService, stationService)

	subscriptionRegistry.Register(stationHandler)
	subscriptionRegistry.Register(rangingHandler)

	mqttClient, err := mqtt.Create(cfg, subscriptionRegistry)
	if err != nil {
		return nil, fmt.Errorf("failed to create MQTT client: %w", err)
	}

	if err := mqttClient.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to MQTT broker: %w", err)
	}

	if err := mqttClient.SubscribeRegistry(); err != nil {
		return nil, fmt.Errorf("failed to subscribe to MQTT topics: %w", err)
	}

	log.Info().Msg("MQTT client initialized and subscribed to topics")
	return mqttClient, nil
}
