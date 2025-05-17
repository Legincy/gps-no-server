package controllers

import (
	"github.com/gin-gonic/gin"
	"gps-no-server/internal/controllers/dto"
	"gps-no-server/internal/services"
	"strconv"
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
		stations.GET("/:id", c.GetById)
	}
}

func (c *RangingController) GetAll(ctx *gin.Context) {
	includeParam := ctx.Query("include")
	sourceParam := ctx.Query("source")
	destinationParam := ctx.Query("destination")

	includes := dto.ParseIncludes(includeParam)

	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully retrieved ranging data",
		"payload": []interface{}{},
	}

	rangingData, err := c.rangingService.GetAll(ctx, includes["stations"], sourceParam, destinationParam)
	if err != nil {
		response["status"] = 500
		response["message"] = err.Error()
		ctx.JSON(500, response)
		return
	}

	transformedResult := dto.FromRangingList(rangingData, includes)

	if len(transformedResult) > 0 {
		response["payload"] = transformedResult
	}

	ctx.JSON(200, response)
}

func (c *RangingController) GetById(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	includeParam := ctx.Query("include")
	includes := dto.ParseIncludes(includeParam)

	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully retrieved ranging data",
	}

	rangingData, err := c.rangingService.GetById(ctx, uint(id))
	if err != nil {
		response["status"] = 500
		response["message"] = err.Error()
		ctx.JSON(500, response)
		return
	}

	transformedResult := dto.FromRanging(rangingData, includes)

	if transformedResult != nil {
		response["payload"] = transformedResult

	}

	ctx.JSON(200, response)
}
