package controllers

import (
	"github.com/gin-gonic/gin"
	"gps-no-server/internal/controllers/dto"
	"gps-no-server/internal/services"
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
	}
}

func (c *ClusterController) GetAll(ctx *gin.Context) {
	includeParam := ctx.Query("include")

	includes := dto.ParseIncludes(includeParam)

	response := map[string]interface{}{
		"status":  200,
		"message": "Successfully retrieved cluster data",
		"payload": []interface{}{},
	}

	clusters, err := c.clusterService.GetAll(ctx, includes["stations"])
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	transformedResult := dto.FromClusterList(clusters, includes)

	if len(transformedResult) > 0 {
		response["payload"] = transformedResult
	}

	ctx.JSON(200, response)
}
