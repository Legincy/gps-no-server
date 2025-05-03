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
		stations.GET("", c.GetAll)
	}
}

func (c *StationController) GetAll(ctx *gin.Context) {
	includeParam := ctx.Query("include")
	includes := dto.ParseIncludes(includeParam)

	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully retrieved station data",
		"payload": []interface{}{},
	}

	stations, err := c.stationService.GetAll(ctx, includes["cluster"])
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	transformedResult := dto.FromStationList(stations, includes)

	if len(transformedResult) > 0 {
		response["payload"] = transformedResult
	}

	ctx.JSON(200, response)

}
