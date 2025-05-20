package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gps-no-server/internal/common/config"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/di"
	"gps-no-server/internal/infrastructure/http/api"
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
	appLog := logger.GetLogger("main")

	container, err := di.NewContainer(cfg)
	if err != nil {
		appLog.Fatal().Err(err).Msg("Failed to initialize application container")
	}
	defer container.Cleanup()

	if err := container.MqttClient.Connect(); err != nil {
		appLog.Error().Err(err).Msg("Error connecting to MQTT broker")
	}
	if err := container.MqttClient.SubscribeRegistry(); err != nil {
		appLog.Error().Err(err).Msg("Error subscribing to MQTT topics")
	}

	server, err := setupServer(cfg, container)
	if err != nil {
		appLog.Fatal().Err(err).Msg("Error while initializing server")
	}

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	appLog.Info().Msg("Shutdown signal received")

	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout*time.Second)
		defer cancel()

		appLog.Info().Msg("Shutting down server...")
		if err := server.Shutdown(ctx); err != nil {
			appLog.Error().Err(err).Msg("Error while shutting down server")
		}
		appLog.Info().Msg("Successfully stopped HTTP server")
	}
}

func setupServer(cfg *config.Config, container *di.Container) (*http.Server, error) {
	gin.SetMode(cfg.Server.ReleaseMode)
	router := gin.New()

	if cfg.Server.LogLevel == "debug" {
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

	apiHandler := api.NewAPI(
		container.StationController,
		container.StationConfigController,
		container.ClusterController,
		container.RangingController,
	)
	apiHandler.RegisterRoutes(router)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout * time.Second,
		WriteTimeout: cfg.Server.WriteTimeout * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info().Msgf("Starting server on %s:%d (%s)", cfg.Server.Host, cfg.Server.Port, cfg.Server.ReleaseMode)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Server failed to start")
		}
	}()

	time.Sleep(100 * time.Millisecond)
	return server, nil
}
