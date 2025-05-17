package dto

import (
	"gorm.io/gorm"
	"gps-no-server/internal/models"
	"time"
)

type StationDto struct {
	gorm.Model `json:"-"`
	ID         uint            `json:"id"`
	MacAddress string          `json:"mac_address" gorm:"unique;not null"`
	Name       string          `json:"name" gorm:"not null"`
	ClusterID  *uint           `json:"cluster_id,omitempty"`
	Cluster    *ClusterDto     `json:"cluster,omitempty"`
	CreatedAt  *time.Time      `json:"created_at,omitempty"`
	UpdatedAt  *time.Time      `json:"updated_at,omitempty"`
	DeletedAt  *gorm.DeletedAt `json:"deleted_at,omitempty"`
	LastSeen   *time.Time      `json:"last_seen,omitempty"`
}

func FromStation(station *models.Station, includes map[string]bool) *StationDto {
	if includes == nil {
		includes = map[string]bool{}
	}

	response := &StationDto{
		ID:         station.ID,
		MacAddress: station.MacAddress,
		Name:       station.Name,
	}

	if includes["cluster"] && station.Cluster != nil {
		response.Cluster = FromCluster(station.Cluster, nil)
	} else {
		response.ClusterID = station.ClusterID
	}

	if includes["meta"] {
		response.CreatedAt = &station.CreatedAt
		response.UpdatedAt = &station.UpdatedAt
		response.DeletedAt = &station.DeletedAt
		response.LastSeen = &station.LastSeen
	}

	return response
}

func ToStation(dto *StationDto) *models.Station {
	return &models.Station{
		MacAddress: dto.MacAddress,
		Name:       dto.Name,
		ClusterID:  dto.ClusterID,
	}
}

func FromStationList(stations []*models.Station, includes map[string]bool) []*StationDto {
	var response []*StationDto

	for _, station := range stations {
		response = append(response, FromStation(station, includes))
	}

	return response
}
