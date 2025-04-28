package repository

import (
	"context"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gps-no-server/internal/logger"
	"gps-no-server/internal/models"
)

type StationRepository struct {
	db  *gorm.DB
	log zerolog.Logger
}

func NewStationRepository(db *gorm.DB) *StationRepository {
	return &StationRepository{
		db:  db,
		log: logger.GetLogger("mqtt"),
	}
}

func (c *StationRepository) Save(ctx context.Context, station *models.Station) error {
	result := c.db.WithContext(ctx).Where("mac_address = ?", station.MacAddress).Assign(station).FirstOrCreate(station)

	return result.Error
}
