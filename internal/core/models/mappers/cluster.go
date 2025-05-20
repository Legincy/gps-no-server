package mappers

import (
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/core/models/dtos"
	"gps-no-server/internal/infrastructure/http/dto"
)

func FromCluster(cluster *models.Cluster, includeParam *string) *dtos.ClusterDto {
	includes := dto.ParseIncludes(includeParam)

	response := &dtos.ClusterDto{
		ID:          cluster.ID,
		Name:        cluster.Name,
		Description: cluster.Description,
	}

	if includes["meta"] {
		response.CreatedAt = &cluster.CreatedAt
		response.UpdatedAt = &cluster.UpdatedAt
		response.DeletedAt = &cluster.DeletedAt
	}

	return response
}

func ToCluster(dto *dtos.ClusterDto) *models.Cluster {
	return &models.Cluster{
		Name:        dto.Name,
		Description: dto.Description,
	}
}

func FromClusterList(clusters []*models.Cluster, includeParam *string) []*dtos.ClusterDto {
	var response []*dtos.ClusterDto

	for _, cluster := range clusters {
		response = append(response, FromCluster(cluster, includeParam))
	}

	return response
}
