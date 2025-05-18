package dtos

import (
	"gorm.io/gorm"
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
