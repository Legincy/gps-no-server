package repositories

import (
	"context"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/core/models"
)

type StationRepository struct {
	*BaseRepository[models.Station]
	db  *gorm.DB
	log zerolog.Logger
}

func NewStationRepository(db *gorm.DB) *StationRepository {
	baseRepository := &BaseRepository[models.Station]{
		DB:         db,
		Log:        logger.GetLogger("station-repository"),
		EntityName: "station-repository",
	}

	return &StationRepository{
		BaseRepository: baseRepository,
		db:             db,
		log:            logger.GetLogger("station-repository"),
	}
}

func (s *StationRepository) FindByMac(ctx context.Context, macAddress string, includes map[string]bool) (*models.Station, error) {
	var station models.Station
	result := s.db.WithContext(ctx).Where("mac_address = ?", macAddress).First(&station)
	return &station, result.Error
}

func (s *StationRepository) FindByIdentifier(ctx context.Context, identifier string, includes map[string]bool) (*models.Station, error) {
	var station models.Station
	result := s.db.WithContext(ctx).Where("identifier = ?", identifier).First(&station)
	return &station, result.Error
}
