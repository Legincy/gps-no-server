package controllers

import (
	"github.com/gin-gonic/gin"
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/core/models/dtos"
	"gps-no-server/internal/core/models/mappers"
	"gps-no-server/internal/core/services"
)

type StationController struct {
	*BaseController[*models.Station, dtos.StationDto]
	stationService *services.StationService
}

func NewStationController(stationService *services.StationService) *StationController {
	baseController := NewBaseController[*models.Station, dtos.StationDto](
		stationService,
		mappers.ToStation,
		mappers.FromStation,
		"/stations",
	)

	return &StationController{
		BaseController: baseController,
		stationService: stationService,
	}
}

func (c *StationController) RegisterRoutes(router *gin.RouterGroup) {
	c.BaseController.RegisterRoutes(router)
	c.Router.GET("/mac/:mac", c.GetByMacAddress)
}

func (c *StationController) GetByMacAddress(ctx *gin.Context) {
	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully retrieved station data",
		"payload": nil,
	}

	macAddress := ctx.Param("mac")
	includeParam := ctx.Query("include")

	station, err := c.stationService.GetByMac(ctx, macAddress, &includeParam)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response["payload"] = station

	ctx.JSON(200, response)
}
