package controllers

import (
	"github.com/gin-gonic/gin"
	"gps-no-server/internal/services"
)

type StationController struct {
	stationService *services.StationService
}

func NewStationController(stationService *services.StationService) *StationController {
	return &StationController{
		stationService: stationService,
	}
}

func (c *StationController) RegisterRoutes(router *gin.RouterGroup) {
	stations := router.Group("/stations")
	{
		stations.GET("", c.GetAllStations)
	}
}

func (c *StationController) GetAllStations(ctx *gin.Context) {
	stations, err := c.stationService.GetAllStations(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, stations)

}
