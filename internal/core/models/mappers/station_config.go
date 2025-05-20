package mappers

import (
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/core/models/dtos"
	"gps-no-server/internal/infrastructure/http/dto"
)

func FromStationConfig(config *models.StationConfiguration, includeParam *string) *dtos.StationConfigurationDto {
	includes := dto.ParseIncludes(includeParam)

	response := &dtos.StationConfigurationDto{
		ID:         config.ID,
		StationID:  config.StationID,
		UWBMode:    string(config.UWBMode),
		UWBChannel: config.UWBChannel,
	}

	if includes["meta"] {
		response.CreatedAt = &config.CreatedAt
		response.UpdatedAt = &config.UpdatedAt
		response.DeletedAt = &config.DeletedAt
	}

	return response
}

func ToStationConfig(dto *dtos.StationConfigurationDto) *models.StationConfiguration {
	return &models.StationConfiguration{
		StationID:  dto.StationID,
		UWBMode:    models.UWBMode(dto.UWBMode),
		UWBChannel: dto.UWBChannel,
	}
}

func FromStationConfigList(configs []*models.StationConfiguration, includeParam *string) []*dtos.StationConfigurationDto {
	response := make([]*dtos.StationConfigurationDto, len(configs))
	for i, config := range configs {
		response[i] = FromStationConfig(config, includeParam)
	}

	return response
}
