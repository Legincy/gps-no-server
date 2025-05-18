package controllers

import (
	"github.com/gin-gonic/gin"
	"gps-no-server/internal/controllers/dto"
	"gps-no-server/internal/services"
	"strconv"
)

type ClusterController struct {
	clusterService *services.ClusterService
}

func NewClusterController(clusterService *services.ClusterService) *ClusterController {
	return &ClusterController{
		clusterService: clusterService,
	}
}

func (c *ClusterController) RegisterRoutes(router *gin.RouterGroup) {
	stations := router.Group("/clusters")
	{
		stations.GET("", c.GetAll)
		stations.GET("/:id", c.GetById)
		stations.POST("", c.Create)
		stations.DELETE("/:id", c.Delete)
		stations.PUT("/:id", c.Update)
	}
}

func (c *ClusterController) GetAll(ctx *gin.Context) {
	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully retrieved cluster data",
		"payload": []interface{}{},
	}

	includeParam := ctx.Query("include")

	clusters, err := c.clusterService.GetAll(ctx, &includeParam)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	transformedResult := dto.FromClusterList(clusters, &includeParam)

	if len(transformedResult) > 0 {
		response["payload"] = transformedResult
	}

	ctx.JSON(200, response)
}

func (c *ClusterController) GetById(ctx *gin.Context) {
	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully retrieved cluster data",
	}

	clusterId := ctx.Param("id")
	if clusterId == "" {
		response["status"] = 400
		response["message"] = "Invalid request payload"
		ctx.JSON(400, response)
		return
	}

	numClusterId, err := strconv.ParseUint(clusterId, 10, 64)

	if err != nil {
		response["status"] = 400
		response["message"] = "Invalid request payload"
		ctx.JSON(400, response)
		return
	}

	includeParam := ctx.Query("include")

	cluster, err := c.clusterService.GetById(ctx, uint(numClusterId), &includeParam)

	if err != nil {
		response["status"] = 404
		response["message"] = "Cluster not found: " + err.Error()
		ctx.JSON(404, response)
		return
	}

	response["payload"] = dto.FromCluster(cluster, &includeParam)
	ctx.JSON(200, response)
}

func (c *ClusterController) Create(ctx *gin.Context) {
	var clusterDto dto.ClusterDto

	response := map[string]interface{}{
		"status":  201,
		"message": "Cluster created successfully",
		"payload": []interface{}{},
	}

	if err := ctx.ShouldBindJSON(&clusterDto); err != nil {
		response["status"] = 400
		response["message"] = "Invalid request payload"
		ctx.JSON(400, response)
		return
	}

	cluster := dto.ToCluster(&clusterDto)
	includeParam := ctx.Query("include")

	cluster, err := c.clusterService.Create(ctx, cluster, &includeParam)

	if err != nil {
		response["status"] = 500
		response["message"] = err.Error()
		ctx.JSON(500, response)
		return
	}

	response["payload"] = dto.FromCluster(cluster, nil)
	ctx.JSON(201, response)
}

func (c *ClusterController) Delete(ctx *gin.Context) {
	response := map[string]interface{}{
		"status":  200,
		"message": "Cluster deleted successfully",
	}

	clusterId := ctx.Param("id")
	if clusterId == "" {
		response["status"] = 400
		response["message"] = "Invalid request payload"
		ctx.JSON(400, response)
		return
	}

	numClusterId, err := strconv.ParseUint(clusterId, 10, 64)

	if err != nil {
		response["status"] = 400
		response["message"] = "Invalid request payload"
		ctx.JSON(400, response)
		return
	}

	includeParam := ctx.Query("include")
	cluster, err := c.clusterService.GetById(ctx, uint(numClusterId), &includeParam)

	if err != nil {
		response["status"] = 500
		response["message"] = err.Error()
		ctx.JSON(500, response)
		return
	}

	err = c.clusterService.Delete(ctx, cluster, &includeParam)

	if err != nil {
		response["status"] = 500
		response["message"] = err.Error()
		ctx.JSON(500, response)
		return
	}

	ctx.JSON(200, response)
}

func (c *ClusterController) Update(ctx *gin.Context) {
	response := map[string]interface{}{
		"status":  200,
		"message": "Cluster updated successfully",
	}

	clusterId := ctx.Param("id")

	if clusterId == "" {
		response["status"] = 400
		response["message"] = "Invalid request payload"
		ctx.JSON(400, response)
		return
	}

	var clusterDto dto.ClusterDto

	if err := ctx.ShouldBindJSON(&clusterDto); err != nil {
		response["status"] = 400
		response["message"] = "Invalid request payload"
		ctx.JSON(400, response)
		return
	}

	numClusterId, err := strconv.ParseUint(clusterId, 10, 64)

	if err != nil {
		response["status"] = 400
		response["message"] = "Invalid request payload"
		ctx.JSON(400, response)
		return
	}

	includeParam := ctx.Query("include")
	cluster, err := c.clusterService.GetById(ctx, uint(numClusterId), &includeParam)
	if err != nil {
		response["status"] = 404
		response["message"] = "Cluster not found: " + err.Error()
		ctx.JSON(404, response)
		return
	}

	if clusterDto.Name != "" {
		cluster.Name = clusterDto.Name
	}

	if clusterDto.Description != "" {
		cluster.Description = clusterDto.Description
	}

	updatedCluster, err := c.clusterService.Update(ctx, cluster, &includeParam)
	if err != nil {
		response["status"] = 500
		response["message"] = "Failed to update cluster: " + err.Error()
		ctx.JSON(500, response)
		return
	}

	response["payload"] = dto.FromCluster(updatedCluster, nil)
	ctx.JSON(200, response)
}
