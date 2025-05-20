package mappers

import (
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/core/models/dtos"
	"gps-no-server/internal/infrastructure/http/dto"
)

func FromStation(station *models.Station, includeParam *string) *dtos.StationDto {
	includes := dto.ParseIncludes(includeParam)

	response := &dtos.StationDto{
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
		response.LastSeen = &station.Uptime
	}

	return response
}

func ToStation(dto *dtos.StationDto) *models.Station {
	return &models.Station{
		MacAddress: dto.MacAddress,
		Name:       dto.Name,
		ClusterID:  dto.ClusterID,
	}
}

func FromStationList(stations []*models.Station, includeParam *string) []*dtos.StationDto {
	var response []*dtos.StationDto

	for _, station := range stations {
		response = append(response, FromStation(station, includeParam))
	}

	return response
}
