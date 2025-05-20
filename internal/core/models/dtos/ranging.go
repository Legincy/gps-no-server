package dtos

import (
	"gorm.io/gorm"
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
