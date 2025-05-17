package dto

import (
	"gorm.io/gorm"
	"gps-no-server/internal/models"
	"time"
)

type ClusterDto struct {
	gorm.Model  `json:"-"`
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Stations    []*StationDto   `json:"stations,omitempty"`
	CreatedAt   *time.Time      `json:"created_at,omitempty"`
	UpdatedAt   *time.Time      `json:"updated_at,omitempty"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty"`
}

func FromCluster(cluster *models.Cluster, includes map[string]bool) *ClusterDto {
	if includes == nil {
		includes = map[string]bool{}
	}

	response := &ClusterDto{
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

func ToCluster(dto *ClusterDto) *models.Cluster {
	return &models.Cluster{
		Name:        dto.Name,
		Description: dto.Description,
	}
}

func FromClusterList(clusters []*models.Cluster, includes map[string]bool) []*ClusterDto {
	var response []*ClusterDto

	for _, cluster := range clusters {
		response = append(response, FromCluster(cluster, includes))
	}

	return response
}
