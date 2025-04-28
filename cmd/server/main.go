package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gps-no-server/internal/config"
	"gps-no-server/internal/database"
	"gps-no-server/internal/logger"
	"gps-no-server/internal/mqtt"
	"gps-no-server/internal/mqtt/subscriptions"
	"gps-no-server/internal/repository"
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

	mqttClient, err := initMqtt(&cfg.Mqtt, gormDB)
	if err != nil {
		log.Error().Msgf("Error while initializing MQTT client: %v", err)
	}
	defer mqttClient.Disconnect()

	server, err := initServer(&cfg.Server)
	if err != nil {
		log.Fatal().Msgf("Error while initializing server: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout*time.Second)
		defer cancel()

		log.Info().Msgf("Shutting down server...")
		if err := server.Shutdown(ctx); err != nil {
			log.Error().Msgf("Error while shutting down server: %v", err)
		}
		log.Info().Msgf("Successfully stopped HTTP server")
	}
}

func initServer(cfg *config.ServerConfig) (*http.Server, error) {
	gin.SetMode(cfg.ReleaseMode)
	router := gin.New()

	if cfg.LogLevel == "debug" {
		router.Use(gin.Logger())
	}

	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "pong",
		})
	})

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout * time.Second,
		WriteTimeout: cfg.WriteTimeout * time.Second,
	}

	errChan := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Msgf("Server failed to start: %v", err)
			errChan <- err
		}
	}()

	log.Info().Msgf("Server started on %s:%d (%s)", cfg.Host, cfg.Port, cfg.ReleaseMode)

	return server, nil
}

func initMqtt(cfg *config.MqttConfig, db *database.GormDB) (*mqtt.Client, error) {
	stationRepository := repository.NewStationRepository(db.DB)

	subscriptionRegistry := mqtt.CreateRegistry()
	subscriptionRegistry.Register(subscriptions.NewStationSubscription(stationRepository))

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

	return mqttClient, nil
}
