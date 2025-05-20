package models

import "gorm.io/gorm"

type UWBMode string

const (
	AnchorMode UWBMode = "ANCHOR"
	TagMode    UWBMode = "TAG"
	NoneMode   UWBMode = "NONE"
)

type StationConfiguration struct {
	gorm.Model
	StationID uint     `gorm:"uniqueIndex;not null"`
	Station   *Station `gorm:"foreignKey:StationID"`

	UWBMode         UWBMode `gorm:"type:varchar(10);not null;default:'ANCHOR'"`
	UWBChannel      uint8   `gorm:"not null;default:5"`
	UWBPreambleCode uint8   `gorm:"not null;default:9"`
	UWBPreambleLen  string  `gorm:"type:varchar(20);not null;default:'128'"`
}

func (s StationConfiguration) SetID(id uint) {
	s.ID = id
}

func (s StationConfiguration) GetID() uint {
	return s.ID
}

func (s StationConfiguration) TableName() string {
	return "station_configurations"
}
