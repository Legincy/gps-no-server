package models

import "gorm.io/gorm"

type Cluster struct {
	gorm.Model
	Name        string    `gorm:"size:100;unique;not null"`
	Description string    `gorm:"type:text"`
	Stations    []Station `gorm:"foreignKey:ClusterID"`
}
