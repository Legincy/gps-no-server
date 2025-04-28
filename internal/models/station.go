package models

import (
	"gorm.io/gorm"
)

type Station struct {
	gorm.Model
	MacAddress string `gorm:"uniqueIndex;not null"`
	Name       string
	PositionX  float64 `gorm:"column:position_x"`
	PositionY  float64 `gorm:"column:position_y"`
}
