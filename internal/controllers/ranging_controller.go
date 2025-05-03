package controllers

import (
	"github.com/gin-gonic/gin"
	"gps-no-server/internal/controllers/dto"
	"gps-no-server/internal/services"
)

type RangingController struct {
	rangingService *services.RangingService
}

func NewRangingController(rangingService *services.RangingService) *RangingController {
	return &RangingController{
		rangingService: rangingService,
	}
}

func (c *RangingController) RegisterRoutes(router *gin.RouterGroup) {
	stations := router.Group("/rangings")
	{
		stations.GET("", c.GetAll)
	}
}

func (c *RangingController) GetAll(ctx *gin.Context) {
	response := make(map[string]interface{})
	includeParams := ctx.Query("include")
	includes := dto.ParseIncludes(includeParams)

	rangings, err := c.rangingService.GetAll(ctx, includes["stations"])
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	transformedResult := dto.FromRangingList(rangings, includes)

	response["payload"] = transformedResult
	if len(transformedResult) == 0 {
		response["payload"] = []interface{}{}
	}
	response["status"] = "success"
	response["message"] = "Clusters retrieved successfully"

	ctx.JSON(200, response)
}
