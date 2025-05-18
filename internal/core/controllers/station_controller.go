package controllers

import (
	"github.com/gin-gonic/gin"
	"gps-no-server/internal/core/models/dtos"
	"gps-no-server/internal/core/models/mappers"
	"gps-no-server/internal/core/services"
	"strconv"
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
		stations.GET("/:id", c.GetById)
		stations.POST("", c.Create)
		stations.DELETE("/:id", c.Delete)
		stations.PUT("/:id", c.Update)
	}
}

func (c *StationController) GetAll(ctx *gin.Context) {
	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully retrieved station data",
		"payload": []interface{}{},
	}

	includeParam := ctx.Query("include")

	stations, err := c.stationService.GetAll(ctx, &includeParam)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	transformedResult := mappers.FromStationList(stations, &includeParam)

	if len(transformedResult) > 0 {
		response["payload"] = transformedResult
	}

	ctx.JSON(200, response)
}

func (c *StationController) GetById(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	includeParam := ctx.Query("include")

	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully retrieved station data",
		"payload": nil,
	}

	station, err := c.stationService.GetById(ctx, uint(id), &includeParam)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	transformedResult := mappers.FromStation(station, &includeParam)

	if transformedResult != nil {
		response["payload"] = transformedResult
	}

	ctx.JSON(200, response)
}

func (c *StationController) Create(ctx *gin.Context) {
	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully created station data",
	}

	var stationDto dtos.StationDto
	if err := ctx.ShouldBindJSON(&stationDto); err != nil {
		response["status"] = 400
		response["message"] = "Invalid request payload"
		ctx.JSON(400, response)
		return
	}

	station := mappers.ToStation(&stationDto)
	includeParam := ctx.Query("include")

	station, err := c.stationService.Create(ctx, station, &includeParam)
	if err != nil {
		response["status"] = 500
		response["message"] = "Failed to create station: " + err.Error()
		ctx.JSON(500, response)
		return
	}

	response["payload"] = mappers.FromStation(station, nil)

	ctx.JSON(200, response)
}

func (c *StationController) Update(ctx *gin.Context) {
	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully updated station data",
	}

	// ID aus dem URL-Parameter extrahieren
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response["status"] = 400
		response["message"] = "Invalid ID format"
		ctx.JSON(400, response)
		return
	}

	var stationDto dtos.StationDto
	if err := ctx.ShouldBindJSON(&stationDto); err != nil {
		response["status"] = 400
		response["message"] = "Invalid request payload: " + err.Error()
		ctx.JSON(400, response)
		return
	}

	includeParam := ctx.Query("include")
	station, err := c.stationService.GetById(ctx, uint(id), &includeParam)
	if err != nil {
		response["status"] = 404
		response["message"] = "Station not found: " + err.Error()
		ctx.JSON(404, response)
		return
	}

	if stationDto.Name != "" {
		station.Name = stationDto.Name
	}

	if stationDto.MacAddress != "" {
		station.MacAddress = stationDto.MacAddress
	}

	if stationDto.ClusterID != nil {
		station.ClusterID = stationDto.ClusterID
	}

	if stationDto.LastSeen != nil {
		station.Uptime = *stationDto.LastSeen
	}

	updatedStation, err := c.stationService.Update(ctx, station, &includeParam)
	if err != nil {
		response["status"] = 500
		response["message"] = "Failed to update station: " + err.Error()
		ctx.JSON(500, response)
		return
	}

	response["payload"] = mappers.FromStation(updatedStation, &includeParam)
	ctx.JSON(200, response)
}

func (c *StationController) Delete(ctx *gin.Context) {
	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully deleted station data",
	}

	idParam := ctx.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	includeParam := ctx.Query("include")
	station, err := c.stationService.GetById(ctx, uint(id), &includeParam)

	if err != nil {
		response["status"] = 404
		response["message"] = "Station not found: " + err.Error()
		ctx.JSON(404, response)
		return
	}

	err = c.stationService.Delete(ctx, station, &includeParam)
	if err != nil {
		response["status"] = 500
		response["message"] = "Failed to delete station: " + err.Error()
		ctx.JSON(500, response)
		return
	}

	ctx.JSON(200, response)
}
