package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gps-no-server/internal/core/interfaces"
	"reflect"
	"strconv"
)

type BaseController[T interfaces.Entity, DTO any] struct {
	Service interfaces.CrudService[T]
	Router  *gin.RouterGroup
	Path    string

	ToEntity   func(*DTO) T
	FromEntity func(T, *string) *DTO
}

func NewBaseController[T interfaces.Entity, DTO any](
	service interfaces.CrudService[T],
	toEntity func(*DTO) T,
	fromEntity func(T, *string) *DTO,
	path string,
) *BaseController[T, DTO] {
	return &BaseController[T, DTO]{
		Service:    service,
		Path:       path,
		ToEntity:   toEntity,
		FromEntity: fromEntity,
	}
}

func (c *BaseController[T, DTO]) RegisterRoutes(router *gin.RouterGroup) {
	c.Router = router.Group(c.Path)
	{
		c.Router.GET("", c.GetAll)
		c.Router.GET("/:id", c.GetById)
		c.Router.POST("", c.Create)
		c.Router.PUT("/:id", c.Update)
		c.Router.DELETE("/:id", c.Delete)
	}
}

func (c *BaseController[T, DTO]) GetAll(ctx *gin.Context) {
	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully retrieved data",
		"payload": []interface{}{},
	}

	includeParam := ctx.Query("include")

	entities, err := c.Service.GetAll(ctx, &includeParam)
	if err != nil {
		response["status"] = 500
		response["message"] = err.Error()
		ctx.JSON(500, response)
		return
	}

	dtos := make([]*DTO, 0, len(entities))
	for _, entity := range entities {
		dtos = append(dtos, c.FromEntity(entity, &includeParam))
	}

	response["payload"] = dtos
	ctx.JSON(200, response)
}

func (c *BaseController[T, DTO]) GetById(ctx *gin.Context) {
	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully retrieved data",
		"payload": nil,
	}

	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response["status"] = 400
		response["message"] = "Invalid ID format"
		ctx.JSON(400, response)
		return
	}

	includeParam := ctx.Query("include")
	entity, err := c.Service.GetById(ctx, uint(id), &includeParam)
	if err != nil {
		response["status"] = 500
		response["message"] = err.Error()
		ctx.JSON(500, response)
		return
	}

	response["payload"] = c.FromEntity(entity, &includeParam)
	ctx.JSON(200, response)
}

func (c *BaseController[T, DTO]) Create(ctx *gin.Context) {
	var dto DTO

	response := map[string]interface{}{
		"status":  201,
		"message": "Successfully created data",
		"payload": nil,
	}

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		response["status"] = 400
		response["message"] = "Invalid request payload: " + err.Error()
		ctx.JSON(400, response)
		return
	}

	entity := c.ToEntity(&dto)
	includeParam := ctx.Query("include")

	createdEntity, err := c.Service.Create(ctx, entity, &includeParam)
	if err != nil {
		response["status"] = 500
		response["message"] = "Failed to create: " + err.Error()
		ctx.JSON(500, response)
		return
	}

	response["payload"] = c.FromEntity(createdEntity, &includeParam)
	ctx.JSON(201, response)
}

func (c *BaseController[T, DTO]) Update(ctx *gin.Context) {
	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully updated data",
		"payload": nil,
	}

	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response["status"] = 400
		response["message"] = "Invalid ID format"
		ctx.JSON(400, response)
		return
	}

	var requestData map[string]interface{}
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		response["status"] = 400
		response["message"] = "Invalid request payload: " + err.Error()
		ctx.JSON(400, response)
		return
	}

	fields := make([]string, 0, len(requestData))
	for field := range requestData {
		fields = append(fields, field)

	}

	var dto DTO
	jsonData, _ := json.Marshal(requestData)
	if err := json.Unmarshal(jsonData, &dto); err != nil {
		response["status"] = 400
		response["message"] = "Invalid request payload: " + err.Error()
		ctx.JSON(400, response)
		return
	}

	entity := c.ToEntity(&dto)

	idField := reflect.ValueOf(entity).Elem().FieldByName("ID")
	if idField.IsValid() && idField.CanSet() {
		idField.SetUint(uint64(id))
	}

	includeParam := ctx.Query("include")

	updatedEntity, err := c.Service.UpdateFields(ctx, entity, fields, &includeParam)
	if err != nil {
		response["status"] = 500
		response["message"] = "Failed to update: " + err.Error()
		ctx.JSON(500, response)
		return
	}

	response["payload"] = c.FromEntity(updatedEntity, &includeParam)
	ctx.JSON(200, response)
}

func (c *BaseController[T, DTO]) Delete(ctx *gin.Context) {
	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully deleted data",
		"payload": nil,
	}

	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response["status"] = 400
		response["message"] = "Invalid ID format"
		ctx.JSON(400, response)
		return
	}

	includeParam := ctx.Query("include")
	entity, err := c.Service.GetById(ctx, uint(id), &includeParam)
	if err != nil {
		response["status"] = 404
		response["message"] = "Entity not found: " + err.Error()
		ctx.JSON(404, response)
		return
	}

	if err := c.Service.Delete(ctx, entity, &includeParam); err != nil {
		response["status"] = 500
		response["message"] = "Failed to delete: " + err.Error()
		ctx.JSON(500, response)
		return
	}

	response["payload"] = entity

	ctx.JSON(200, response)
}
