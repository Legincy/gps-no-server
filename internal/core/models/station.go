package models

import (
	"gorm.io/gorm"
	"time"
)

type Station struct {
	gorm.Model
	MacAddress    string `gorm:"uniqueIndex;not null"`
	Name          string `gorm:"size:100;not null"`
	ClusterID     *uint
	Cluster       *Cluster `gorm:"foreignKey:ClusterID"`
	Uptime        time.Time
	StationConfig *StationConfiguration `gorm:"foreignKey:StationID"`
}

func (s Station) SetID(id uint) {
	s.ID = id
}

func (s Station) GetID() uint {
	return s.ID
}

func (s Station) TableName() string {
	return "stations"
}

func (s *Station) AfterCreate(tx *gorm.DB) error {
	var count int64
	if err := tx.Model(&StationConfiguration{}).Where("station_id = ?", s.ID).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		config := StationConfiguration{
			StationID:       s.ID,
			UWBMode:         AnchorMode,
			UWBChannel:      5,
			UWBPreambleCode: 9,
			UWBPreambleLen:  "128",
		}

		return tx.Create(&config).Error
	}

	return nil
}
