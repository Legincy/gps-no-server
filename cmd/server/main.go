package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gps-no-server/internal/common/config"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/di"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	logger.Init(cfg.Server.LogLevel)
	appLog := logger.GetLogger("main")

	container, err := di.NewContainer(cfg)
	if err != nil {
		appLog.Fatal().Err(err).Msg("Failed to initialize application container")
	}
	defer func() {
		if err := container.Cleanup(); err != nil {
			appLog.Error().Err(err).Msg("Error during cleanup")
		}
	}()

	if err := container.GetMQTTConnection().Connect(); err != nil {
		appLog.Error().Err(err).Msg("Failed to connect to MQTT broker")
	}

	server, err := setupServer(cfg, container)
	if err != nil {
		appLog.Fatal().Err(err).Msg("Failed to setup server")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	appLog.Info().Msg("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLog.Error().Err(err).Msg("Server forced to shutdown")
	}

	appLog.Info().Msg("Server exited")
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

	api := router.Group("/api/v1")
	{
		container.GetStationHandler().RegisterRoutes(api)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout * time.Second,
		WriteTimeout: cfg.Server.WriteTimeout * time.Second,
		IdleTimeout:  cfg.Server.IdleTimeout * time.Second,
	}

	go func() {
		log.Info().Msgf("Starting server on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Server failed to start")
		}
	}()

	time.Sleep(100 * time.Millisecond)
	return server, nil
}
