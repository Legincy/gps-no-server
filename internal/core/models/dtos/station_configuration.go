package dtos

import (
	"gorm.io/gorm"
	"time"
)

type StationConfigurationDto struct {
	gorm.Model `json:"-"`
	ID         uint            `json:"id"`
	StationID  uint            `json:"station_id"`
	UWBMode    string          `json:"uwb_mode"`
	UWBChannel uint8           `json:"uwb_channel"`
	CreatedAt  *time.Time      `json:"created_at,omitempty"`
	UpdatedAt  *time.Time      `json:"updated_at,omitempty"`
	DeletedAt  *gorm.DeletedAt `json:"deleted_at,omitempty"`
}
