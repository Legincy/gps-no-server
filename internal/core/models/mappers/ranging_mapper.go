package mappers

import (
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/core/models/dtos"
	"gps-no-server/internal/infrastructure/http/dto"
)

func FromRanging(ranging *models.Ranging, includeParam *string) *dtos.RangingDto {
	includes := dto.ParseIncludes(includeParam)

	response := &dtos.RangingDto{
		ID:          ranging.ID,
		RawDistance: ranging.RawDistance,
	}

	if includes["stations"] {
		if ranging.Source != nil {
			response.Source = FromStation(ranging.Source, nil)
		}

		if ranging.Destination != nil {
			response.Destination = FromStation(ranging.Destination, nil)
		}
	} else {
		response.SourceID = ranging.SourceID
		response.DestinationID = ranging.DestinationID
	}

	if includes["meta"] {
		response.CreatedAt = &ranging.CreatedAt
		response.UpdatedAt = &ranging.UpdatedAt
		response.DeletedAt = &ranging.DeletedAt
	}

	return response
}

func FromRangingList(rangingList []*models.Ranging, includeParam *string) []*dtos.RangingDto {
	var response []*dtos.RangingDto

	for _, ranging := range rangingList {
		response = append(response, FromRanging(ranging, includeParam))
	}

	return response
}
