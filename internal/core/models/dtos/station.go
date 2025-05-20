package dtos

import (
	"gorm.io/gorm"
	"time"
)

type StationDto struct {
	gorm.Model `json:"-"`
	ID         uint            `json:"id"`
	MacAddress string          `json:"mac_address" gorm:"unique;not null"`
	Name       string          `json:"name" gorm:"not null"`
	ClusterID  *uint           `json:"cluster_id,omitempty"`
	Cluster    *ClusterDto     `json:"cluster,omitempty"`
	CreatedAt  *time.Time      `json:"created_at,omitempty"`
	UpdatedAt  *time.Time      `json:"updated_at,omitempty"`
	DeletedAt  *gorm.DeletedAt `json:"deleted_at,omitempty"`
	LastSeen   *time.Time      `json:"last_seen,omitempty"`
}
