package controllers

import (
	"github.com/gin-gonic/gin"
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/core/models/dtos"
	"gps-no-server/internal/core/models/mappers"
	"gps-no-server/internal/core/services"
	"strconv"
)

type RangingController struct {
	*BaseController[*models.Ranging, dtos.RangingDto]
	rangingService *services.RangingService
	eventService   *services.EventStreamService
}

func NewRangingController(rangingService *services.RangingService, eventService *services.EventStreamService) *RangingController {
	baseController := NewBaseController[*models.Ranging, dtos.RangingDto](
		rangingService,
		mappers.ToRanging,
		mappers.FromRanging,
		"/rangings",
	)

	return &RangingController{
		BaseController: baseController,
		rangingService: rangingService,
		eventService:   eventService,
	}
}

func (c *RangingController) RegisterRoutes(router *gin.RouterGroup) {
	c.BaseController.RegisterRoutes(router)

	api := router.Group("/rangings")
	{
		api.GET("/stream", c.StreamAllRangingEvents)
		api.GET("/stream/:id", c.StreamRangingById)
	}
}

func (c *RangingController) StreamAllRangingEvents(ctx *gin.Context) {
	c.eventService.HandleSSERequest(ctx, services.RangingEventType)
}

func (c *RangingController) StreamRangingById(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	c.eventService.HandleSSERequest(ctx, services.RangingEventType, uint(id))
}
