// @title GPS-NO Server API
// @version 1.0
// @description Dies ist die API für den GPS-NO Server.
// @host localhost:8080
// @BasePath /api/v1

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	//Swagger
	//_ "gps-no-server/internal/infrastructure/http/api/docs"
)

type RouterRegistry interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type API struct {
	registry []RouterRegistry
}

func NewAPI(registry ...RouterRegistry) *API {
	return &API{
		registry,
	}
}

func (api *API) RegisterRoutes(router *gin.Engine) {
	apiGroup := router.Group("/api/v1")

	for _, registry := range api.registry {
		registry.RegisterRoutes(apiGroup)
	}
}

// ----------------------------------------------------------------------------- Station Routen -----------------------------------------------------------------------------

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

// Bei routen änderungen vorher swag init -g internal/infrastructure/http/api/router.go ausfüren damit die Anzeige aktualisiert wird
