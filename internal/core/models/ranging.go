package models

import (
	"gorm.io/gorm"
)

type Ranging struct {
	gorm.Model
	SourceID      *uint    `gorm:"not null"`
	Source        *Station `gorm:"foreignKey:SourceID"`
	DestinationID *uint    `gorm:"not null"`
	Destination   *Station `gorm:"foreignKey:DestinationID"`
	RawDistance   float64  `gorm:"default:0.0"`
}

func (r Ranging) SetID(id uint) {
	r.ID = id
}

func (r Ranging) GetID() uint {
	return r.ID
}

func (r Ranging) TableName() string {
	return "rangings"
}
