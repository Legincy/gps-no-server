package controllers

import (
	"github.com/gin-gonic/gin"
	"gps-no-server/internal/controllers/dto"
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
	response := make(map[string]interface{})
	includeParams := ctx.Query("include")
	includes := dto.ParseIncludes(includeParams)

	stations, err := c.stationService.GetAllStations(ctx, includes["cluster"])
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	transformedResult := dto.FromStationList(stations, includes)

	response["payload"] = transformedResult
	if len(transformedResult) == 0 {
		response["payload"] = []interface{}{}
	}
	response["status"] = "success"
	response["message"] = "Stations retrieved successfully"

	ctx.JSON(200, response)

}
