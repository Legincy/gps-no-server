package controllers

import (
	"github.com/gin-gonic/gin"
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/core/models/dtos"
	"gps-no-server/internal/core/models/mappers"
	"gps-no-server/internal/core/services"
)

type StationConfigController struct {
	*BaseController[*models.StationConfiguration, dtos.StationConfigurationDto]
	StationConfigService *services.StationConfigurationService
}

func NewStationConfigController(stationConfigService *services.StationConfigurationService) *StationConfigController {
	baseController := NewBaseController[*models.StationConfiguration, dtos.StationConfigurationDto](
		stationConfigService,
		mappers.ToStationConfig,
		mappers.FromStationConfig,
		"/station-configurations",
	)

	return &StationConfigController{
		BaseController:       baseController,
		StationConfigService: stationConfigService,
	}
}

func (c *StationConfigController) RegisterRoutes(router *gin.RouterGroup) {
	c.BaseController.RegisterRoutes(router)
}
