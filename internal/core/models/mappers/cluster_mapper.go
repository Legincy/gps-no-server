package mappers

import (
	"gps-no-server/internal/core/models"
	dto2 "gps-no-server/internal/core/models/dtos"
	"gps-no-server/internal/infrastructure/http/dto"
)

func FromCluster(cluster *models.Cluster, includeParam *string) *dto2.ClusterDto {
	includes := dto.ParseIncludes(includeParam)

	response := &dto2.ClusterDto{
		ID:          cluster.ID,
		Name:        cluster.Name,
		Description: cluster.Description,
	}

	if includes["stations"] {

		stationPointers := make([]*models.Station, len(cluster.Stations))
		for i := range cluster.Stations {
			stationPointers[i] = &cluster.Stations[i]
		}

		response.Stations = FromStationList(stationPointers, nil)
	}

	if includes["meta"] {
		response.CreatedAt = &cluster.CreatedAt
		response.UpdatedAt = &cluster.UpdatedAt
		response.DeletedAt = &cluster.DeletedAt
	}

	return response
}

func ToCluster(dto *dto2.ClusterDto) *models.Cluster {
	return &models.Cluster{
		Name:        dto.Name,
		Description: dto.Description,
	}
}

func FromClusterList(clusters []*models.Cluster, includeParam *string) []*dto2.ClusterDto {
	var response []*dto2.ClusterDto

	for _, cluster := range clusters {
		response = append(response, FromCluster(cluster, includeParam))
	}

	return response
}
