package models

import (
	"gorm.io/gorm"
)

type Cluster struct {
	gorm.Model
	Name        string    `gorm:"size:100;unique;not null"`
	Description string    `gorm:"type:text"`
	Stations    []Station `gorm:"foreignKey:ClusterID"`
}

func (c Cluster) SetID(id uint) {
	c.ID = id
}

func (c Cluster) GetID() uint {
	return c.ID
}

func (c Cluster) TableName() string {
	return "clusters"
}
