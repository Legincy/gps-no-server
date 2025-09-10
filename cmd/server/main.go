package main

import (
	"context"
	"fmt"
	"gps-no-server/internal/common/config"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/di"
	"gps-no-server/internal/infrastructure/http/api"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	_ "gps-no-server/docs" // <- Wichtig: Pfad anpassen zu deinem Modulnamen
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

	r := gin.Default()
	api := r.Group("/api/v1")
	{
		api.POST("stations", CreateStation)
		api.GET("stations", GetStations)
		api.GET("stations/:mac", GetStation)
		api.PUT("stations/:mac", UpdateStation)
		api.PUT("stations/:mac", DeleteStation)

		api.POST("/clusters", CreateCluster)
		api.GET("/clusters/:id", GetCluster)
		api.GET("/clusters", GetClusters)
		api.PUT("/clusters/:id", UpdateCluster)
		api.DELETE("/clusters/:id", DeleteCluster)

		api.GET("measurements/stations/:mac", GetMeasurement)
		api.GET("measurements/cluster/:mac", GetMeasurementCluster)
	}
	r.StaticFile("/swagger/doc.json", "./docs/swagger.json")

	r.Run(":8090")

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

// @Summary	Erstellt eine Station
// @Description	Erstellt eine Station
// @Tags	Station
// @Produce	json
// @Success	200 {string} string "Station erstellt"
// @Router 	/stations [post]
func CreateStation(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Station erstellt"})
}

// @Summary	Gibt alle Stationen wieder
// @Description	Gibt alle Stationen wieder
// @Tags	Station
// @Produce	json
// @Success	200 {string} string "Alle Stationen wiedergegeben"
// @Router 	/stations [get]
func GetStations(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Stationen zurückgegeben"})
}

// @Summary	Gibt eine Station wieder
// @Description	Gibt eine Station basierend auf der angegebenen Mac Adresse wieder
// @Tags	Station
// @Produce	json
// @Success	200 {string} string "Station wiedergegeben"
// @Router 	/stations/{Mac} [get]
func GetStation(c *gin.Context) {
	mac := c.Param("mac")
	c.JSON(http.StatusOK, gin.H{"message": "Station zurrückgegeben", "mac": mac})
}

// @Summary	Updated eine Station
// @Description	Updated eine Station basierend auf der angegebenen Mac Adresse
// @Tags	Station
// @Produce	json
// @Success	200 {string} string "Station geupdated"
// @Router 	/stations/{Mac} [put]
func UpdateStation(c *gin.Context) {
	mac := c.Param("mac")
	c.JSON(http.StatusOK, gin.H{"message": "Station geupdated", "mac": mac})
}

// @Summary	Löscht eine Station
// @Description	Löscht eine Station basierend auf der angegebenen Mac Adresse
// @Tags	Station
// @Produce	json
// @Success	200 {string} string "Station gelöscht"
// @Router 	/stations/{Mac} [delete]
func DeleteStation(c *gin.Context) {
	mac := c.Param("mac")
	c.JSON(http.StatusOK, gin.H{"message": "Station gelöscht", "mac": mac})
}

// ----------------------------------------------------------------------------- Cluster Routen -----------------------------------------------------------------------------

// @Summary	Erstellt ein Cluster
// @Description	Erstellt ein Cluster
// @Tags	Cluster
// @Produce	json
// @Success	200 {string} string "Cluster erstellt"
// @Router 	/clusters [post]
func CreateCluster(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Cluster erstellt"})
}

// @Summary	Gibt ein Cluster wieder
// @Description	Gibt ein Cluster basierend auf der angegebenen ID wieder
// @Tags	Cluster
// @Produce	json
// @Success	200 {string} string "Cluster wiedergegeben"
// @Router 	/clusters/{ID} [get]
func GetCluster(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Cluster wiedergegeben", "id": id})
}

// @Summary	Gibt alle Cluster wieder
// @Description	Gibt alle Cluster wieder
// @Tags	Cluster
// @Produce	json
// @Success	200 {string} string "Cluster wiedergegeben"
// @Router 	/clusters [get]
func GetClusters(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Cluster wiedergegeben"})
}

// @Summary	Updated ein Cluster
// @Description	Updated ein Cluster basierend auf der angegebenen ID
// @Tags	Cluster
// @Produce	json
// @Success	200 {string} string "Cluster geupdated"
// @Router 	/clusters/{ID} [put]
func UpdateCluster(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Cluster geupdated", "id": id})
}

// @Summary	Löscht ein Cluster
// @Description	Löscht ein Cluster basierend auf der angegebenen ID
// @Tags	Cluster
// @Produce	json
// @Success	200 {string} string "Cluster gelöscht"
// @Router 	/clusters/{ID} [delete]
func DeleteCluster(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Cluster gelöscht", "id": id})
}

// ----------------------------------------------------------------------------- Measurement Routen -----------------------------------------------------------------------------

// @Summary	Gibt die Messung einer Station wieder
// @Description	Gibt die Messung einer Station basierend auf der angegebenen Mac Adresse wieder
// @Tags	Station, Measurement
// @Produce	json
// @Success	200 {string} string "Messung wiedergegeben"
// @Router 	/measurements/stations/{Mac} [get]
func GetMeasurement(c *gin.Context) {
	mac := c.Param("mac")
	c.JSON(http.StatusOK, gin.H{"message": "Messung zurückgegeben", "mac": mac})
}

// @Summary	Gibt die Messungen eines Clusters wieder
// @Description	Gibt die Messungen eines Clusters basierend auf der angegebenen Mac Adresse wieder
// @Tags	Station, Measurement, Cluster
// @Produce	json
// @Success	200 {string} string "Messungen wiedergegeben"
// @Router 	/measurements/cluster/{Mac} [get]
func GetMeasurementCluster(c *gin.Context) {
	mac := c.Param("mac")
	c.JSON(http.StatusOK, gin.H{"message": "Messungen zurückgegeben", "mac": mac})
}
