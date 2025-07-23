package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/core/models/dtos"
	"gps-no-server/internal/core/models/mappers"
	"gps-no-server/internal/di/interfaces"
	"strconv"
)

type StationHandler struct {
	service interfaces.StationService
	log     zerolog.Logger
}

func NewStationHandler(service interfaces.StationService) *StationHandler {
	return &StationHandler{
		service: service,
		log:     logger.GetLogger("station-handler"),
	}
}

func (h *StationHandler) RegisterRoutes(router *gin.RouterGroup) {
	stations := router.Group("/stations")
	{
		stations.GET("", h.GetAll)
		stations.GET("/:id", h.GetByID)
		stations.GET("/mac/:mac", h.GetByMac)
		stations.POST("", h.Create)
		stations.PUT("/:id", h.Update)
		stations.DELETE("/:id", h.Delete)
	}
}

func (h *StationHandler) GetAll(c *gin.Context) {
	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully retrieved stations",
		"data":    []interface{}{},
	}

	includeParam := c.Query("include")

	stations, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to retrieve stations")
		response["status"] = 500
		response["message"] = "Failed to retrieve stations: " + err.Error()
		c.JSON(500, response)
		return
	}

	stationDTOs := make([]*dtos.StationDto, 0, len(stations))
	for _, station := range stations {
		dto := mappers.FromStation(station, &includeParam)
		stationDTOs = append(stationDTOs, dto)
	}

	response["data"] = stationDTOs
	c.JSON(200, response)
}

func (h *StationHandler) GetByID(c *gin.Context) {
	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully retrieved station",
		"data":    nil,
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.log.Warn().Str("id", idParam).Msg("Invalid ID format")
		response["status"] = 400
		response["message"] = "Invalid ID format"
		c.JSON(400, response)
		return
	}

	includeParam := c.Query("include")

	station, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.log.Error().Err(err).Uint64("id", id).Msg("Failed to retrieve station")
		response["status"] = 404
		response["message"] = "Station not found"
		c.JSON(404, response)
		return
	}

	dto := mappers.FromStation(station, &includeParam)
	response["data"] = dto
	c.JSON(200, response)
}

func (h *StationHandler) GetByMac(c *gin.Context) {
	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully retrieved station",
		"data":    nil,
	}

	macAddress := c.Param("mac")
	if macAddress == "" {
		response["status"] = 400
		response["message"] = "MAC address is required"
		c.JSON(400, response)
		return
	}

	includeParam := c.Query("include")

	station, err := h.service.GetByMac(c.Request.Context(), macAddress)
	if err != nil {
		h.log.Error().Err(err).Str("mac", macAddress).Msg("Failed to retrieve station by MAC")
		response["status"] = 404
		response["message"] = "Station not found"
		c.JSON(404, response)
		return
	}

	dto := mappers.FromStation(station, &includeParam)
	response["data"] = dto
	c.JSON(200, response)
}

func (h *StationHandler) Create(c *gin.Context) {
	response := map[string]interface{}{
		"status":  201,
		"message": "Successfully created station",
		"data":    nil,
	}

	var stationDTO dtos.StationDto
	if err := c.ShouldBindJSON(&stationDTO); err != nil {
		h.log.Warn().Err(err).Msg("Invalid request payload")
		response["status"] = 400
		response["message"] = "Invalid request payload: " + err.Error()
		c.JSON(400, response)
		return
	}

	if stationDTO.MacAddress == "" {
		response["status"] = 400
		response["message"] = "MAC address is required"
		c.JSON(400, response)
		return
	}

	if stationDTO.Name == "" {
		response["status"] = 400
		response["message"] = "Station name is required"
		c.JSON(400, response)
		return
	}

	station := mappers.ToStation(&stationDTO)

	includeParam := c.Query("include")

	createdStation, err := h.service.Create(c.Request.Context(), station)
	if err != nil {
		h.log.Error().Err(err).Str("mac", station.MacAddress).Msg("Failed to create station")
		response["status"] = 500
		response["message"] = "Failed to create station: " + err.Error()
		c.JSON(500, response)
		return
	}

	dto := mappers.FromStation(createdStation, &includeParam)
	response["data"] = dto
	c.JSON(201, response)
}

func (h *StationHandler) Update(c *gin.Context) {
	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully updated station",
		"data":    nil,
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response["status"] = 400
		response["message"] = "Invalid ID format"
		c.JSON(400, response)
		return
	}

	var stationDTO dtos.StationDto
	if err := c.ShouldBindJSON(&stationDTO); err != nil {
		response["status"] = 400
		response["message"] = "Invalid request payload: " + err.Error()
		c.JSON(400, response)
		return
	}

	station := mappers.ToStation(&stationDTO)
	station.ID = uint(id)

	includeParam := c.Query("include")

	updatedStation, err := h.service.Update(c.Request.Context(), station)
	if err != nil {
		h.log.Error().Err(err).Uint64("id", id).Msg("Failed to update station")
		response["status"] = 500
		response["message"] = "Failed to update station: " + err.Error()
		c.JSON(500, response)
		return
	}

	dto := mappers.FromStation(updatedStation, &includeParam)
	response["data"] = dto
	c.JSON(200, response)
}

func (h *StationHandler) Delete(c *gin.Context) {
	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully deleted station",
		"data":    nil,
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response["status"] = 400
		response["message"] = "Invalid ID format"
		c.JSON(400, response)
		return
	}

	station, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response["status"] = 404
		response["message"] = "Station not found"
		c.JSON(404, response)
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		h.log.Error().Err(err).Uint64("id", id).Msg("Failed to delete station")
		response["status"] = 500
		response["message"] = "Failed to delete station: " + err.Error()
		c.JSON(500, response)
		return
	}

	includeParam := c.Query("include")
	dto := mappers.FromStation(station, &includeParam)
	response["data"] = dto
	c.JSON(200, response)
}
