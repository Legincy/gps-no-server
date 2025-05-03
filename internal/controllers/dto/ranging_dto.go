package dto

import (
	"gorm.io/gorm"
	"gps-no-server/internal/models"
	"time"
)

type RangingDto struct {
	ID            uint            `json:"id"`
	SourceID      *uint           `json:"source_id,omitempty"`
	Source        *StationDto     `json:"source,omitempty"`
	DestinationID *uint           `json:"destination_id,omitempty"`
	Destination   *StationDto     `json:"destination,omitempty"`
	RawDistance   float64         `json:"raw_distance" validate:"required"`
	CreatedAt     *time.Time      `json:"created_at,omitempty"`
	UpdatedAt     *time.Time      `json:"updated_at,omitempty"`
	DeletedAt     *gorm.DeletedAt `json:"deleted_at,omitempty"`
}

func FromRanging(ranging *models.Ranging, includes map[string]bool) *RangingDto {
	if includes == nil {
		includes = map[string]bool{}
	}

	response := &RangingDto{
		RawDistance: ranging.RawDistance,
	}

	if includes["stations"] {
		response.Source = FromStation(ranging.Source, nil)
		response.Destination = FromStation(ranging.Destination, nil)
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

func FromRangingList(rangings []*models.Ranging, includes map[string]bool) []*RangingDto {
	var response []*RangingDto

	for _, ranging := range rangings {
		response = append(response, FromRanging(ranging, includes))
	}

	return response
}
