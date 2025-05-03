package dto

import (
	"gorm.io/gorm"
	"gps-no-server/internal/models"
)

type ClusterDto struct {
	gorm.Model  `json:"-"`
	ID          uint          `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Stations    []*StationDto `json:"stations,omitempty"`
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

	return response
}

func FromClusterList(clusters []*models.Cluster, includes map[string]bool) []*ClusterDto {
	var response []*ClusterDto

	for _, cluster := range clusters {
		response = append(response, FromCluster(cluster, includes))
	}

	return response
}
