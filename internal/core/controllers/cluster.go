package controllers

import (
	"github.com/gin-gonic/gin"
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/core/models/dtos"
	"gps-no-server/internal/core/models/mappers"
	"gps-no-server/internal/core/services"
)

type ClusterController struct {
	*BaseController[*models.Cluster, dtos.ClusterDto]
	clusterService *services.ClusterService
}

func NewClusterController(clusterService *services.ClusterService) *ClusterController {
	baseController := NewBaseController[*models.Cluster, dtos.ClusterDto](
		clusterService,
		mappers.ToCluster,
		mappers.FromCluster,
		"/clusters",
	)

	return &ClusterController{
		BaseController: baseController,
		clusterService: clusterService,
	}
}

func (c *ClusterController) RegisterRoutes(router *gin.RouterGroup) {
	c.BaseController.RegisterRoutes(router)
}
