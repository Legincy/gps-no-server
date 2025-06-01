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

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		if err := container.MqttClient.Connect(); err != nil {
			appLog.Warn().Err(err).Msgf("MQTT connection attempt %d/%d failed", i+1, maxRetries)
			if i == maxRetries-1 {
				appLog.Error().Err(err).Msg("Failed to connect to MQTT after max retries")
			} else {
				time.Sleep(time.Duration(i+1) * time.Second)
			}
		} else {
			appLog.Info().Msg("MQTT connected successfully")
			break
		}
	}

	if err := container.MqttClient.SubscribeRegistry(); err != nil {
		appLog.Error().Err(err).Msg("Error subscribing to MQTT topics")
	}

	server, err := setupServer(cfg, container)
	if err != nil {
		appLog.Fatal().Err(err).Msg("Error while initializing server")
	}

	setupHealthCheck(server, container)

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

	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Error().Str("error", err).Msg("Panic recovered")
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Internal server error",
		})
	}))

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

	router.Use(func(c *gin.Context) {
		if c.Request.URL.Path == "/api/v1/rangings/stream" ||
			c.Request.URL.Path == "/api/v1/rangings/stream/*" {
			c.Next()
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
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
		Addr:              fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:           router,
		ReadTimeout:       cfg.Server.ReadTimeout * time.Second,
		WriteTimeout:      cfg.Server.WriteTimeout * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20,
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

func setupHealthCheck(server *http.Server, container *di.Container) {
	if ginRouter, ok := server.Handler.(*gin.Engine); ok {
		ginRouter.GET("/health", func(c *gin.Context) {
			health := gin.H{
				"status":    "ok",
				"timestamp": time.Now().UTC(),
			}

			if container.EventStreamService != nil {
				health["event_stream"] = container.EventStreamService.GetStats()
			}

			if container.Database != nil {
				if db, err := container.Database.DB.DB(); err == nil {
					if err := db.Ping(); err != nil {
						health["database"] = "error"
						health["database_error"] = err.Error()
						c.JSON(503, health)
						return
					}
					health["database"] = "ok"
				}
			}

			if container.MqttClient != nil {
				health["mqtt"] = "ok"
			}

			c.JSON(200, health)
		})
	}
}
