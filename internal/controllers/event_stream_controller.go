package controllers

import (
	"github.com/gin-gonic/gin"
	"gps-no-server/internal/services"
	"strconv"
)

type EventStreamController struct {
	eventService *services.EventStreamService
}

func NewEventStreamController(eventService *services.EventStreamService) *EventStreamController {
	return &EventStreamController{
		eventService: eventService,
	}
}

func (c *EventStreamController) RegisterRoutes(router *gin.RouterGroup) {
	rangings := router.Group("/rangings")
	{
		rangings.GET("/stream", c.StreamAllRangingEvents)
		rangings.GET("/:id/stream", c.StreamRangingById)
	}
}

func (c *EventStreamController) StreamAllRangingEvents(ctx *gin.Context) {
	c.eventService.HandleSSERequest(ctx, services.RangingEventType)
}

func (c *EventStreamController) StreamRangingById(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	c.eventService.HandleSSERequest(ctx, services.RangingEventType, uint(id))
}
