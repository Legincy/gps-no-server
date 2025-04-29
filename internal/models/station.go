package models

import (
	"gorm.io/gorm"
	"time"
)

type Station struct {
	gorm.Model
	MacAddress string `gorm:"uniqueIndex;not null"`
	Name       string
	ClusterID  *uint
	Cluster    *Cluster `gorm:"foreignKey:ClusterID"`
	LastSeen   time.Time
}
