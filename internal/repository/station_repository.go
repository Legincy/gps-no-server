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
		log: logger.GetLogger("station-repository"),
	}
}

func (c *StationRepository) Save(ctx context.Context, station *models.Station) error {
	result := c.db.WithContext(ctx).Where("mac_address = ?", station.MacAddress).Assign(station).FirstOrCreate(station)

	return result.Error
}

func (c *StationRepository) FindByMacAddress(ctx context.Context, macAddress string) (*models.Station, error) {
	station := &models.Station{}
	result := c.db.WithContext(ctx).Where("mac_address = ?", macAddress).First(station)

	if result.Error != nil {
		return nil, result.Error
	}

	return station, nil
}
