package controllers

import (
	"github.com/gin-gonic/gin"
	"gps-no-server/internal/services"
)

type EventStreamController struct {
	rangingService *services.RangingService
}

func NewEventStreamController(rangingService *services.RangingService) *EventStreamController {
	return &EventStreamController{
		rangingService: rangingService,
	}
}

func (c *EventStreamController) RegisterRoutes(router *gin.RouterGroup) {
	//stream := router.Group("/stream")
	{
		//stream.GET("/ranging", c.GetRangingStream)
	}
}

func (c *EventStreamController) StreamRangingEvents(ctx *gin.Context) {
}
