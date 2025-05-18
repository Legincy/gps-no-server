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

func (s *StationRepository) Save(ctx context.Context, station *models.Station) (*models.Station, error) {
	result := s.db.WithContext(ctx).Where("mac_address = ?", station.MacAddress).FirstOrCreate(station)

	return station, result.Error
}

func (s *StationRepository) FindAll(ctx context.Context, includes map[string]bool) ([]*models.Station, error) {
	var stations []*models.Station

	query := s.db.WithContext(ctx)

	if includes["cluster"] {
		query = query.Preload("Cluster")
	}

	result := query.Find(&stations)

	return stations, result.Error
}

func (s *StationRepository) FindByID(ctx context.Context, id uint, includes map[string]bool) (*models.Station, error) {
	var station models.Station
	result := s.db.WithContext(ctx).Where("id = ?", id).First(&station)

	if result.Error != nil {
		return nil, result.Error
	}

	return &station, result.Error
}

func (s *StationRepository) FindByMac(ctx context.Context, macAddress string, includes map[string]bool) (*models.Station, error) {
	var station models.Station
	result := s.db.WithContext(ctx).Where("mac_address = ?", macAddress).First(&station)

	return &station, result.Error
}

func (s *StationRepository) FindActive(ctx context.Context) ([]*models.Station, error) {
	var stations []*models.Station
	result := s.db.WithContext(ctx).Where("active = ?", true).Find(&stations)

	return stations, result.Error
}

func (c *StationRepository) FindByMacAddress(ctx context.Context, macAddress string) (*models.Station, error) {
	station := &models.Station{}
	result := c.db.WithContext(ctx).Where("mac_address = ?", macAddress).First(station)

	if result.Error != nil {
		return nil, result.Error
	}

	return station, nil
}

func (c *StationRepository) DeleteById(ctx context.Context, id uint) error {
	result := c.db.WithContext(ctx).Delete(&models.Station{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (c *StationRepository) Update(ctx context.Context, station *models.Station) (*models.Station, error) {
	result := c.db.WithContext(ctx).Save(station)

	if result.Error != nil {
		return nil, result.Error
	}

	return station, nil
}

func (c *StationRepository) Create(ctx context.Context, station *models.Station) (*models.Station, error) {
	result := c.db.WithContext(ctx).Create(station)

	if result.Error != nil {
		return nil, result.Error
	}

	return station, nil
}
