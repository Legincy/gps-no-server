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
		stations.GET("", c.GetAllClusters)
	}
}

func (c *ClusterController) GetAllClusters(ctx *gin.Context) {
	response := make(map[string]interface{})
	includeParams := ctx.Query("include")
	includes := dto.ParseIncludes(includeParams)

	clusters, err := c.clusterService.GetAllClusters(ctx, includes["stations"])
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	transformedResult := dto.FromClusterList(clusters, includes)

	response["payload"] = transformedResult
	if len(transformedResult) == 0 {
		response["payload"] = []interface{}{}
	}
	response["status"] = "success"
	response["message"] = "Clusters retrieved successfully"

	ctx.JSON(200, response)
}
