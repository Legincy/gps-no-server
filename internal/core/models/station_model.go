package models

import (
	"gorm.io/gorm"
	"time"
)

type Station struct {
	gorm.Model
	MacAddress string `gorm:"uniqueIndex;not null"`
	Name       string `gorm:"size:100;not null"`
	ClusterID  *uint
	Cluster    *Cluster `gorm:"foreignKey:ClusterID"`
	Uptime     time.Time
}
