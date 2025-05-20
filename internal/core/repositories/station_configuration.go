package repositories

import (
	"context"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/core/models"
)

type StationConfigurationRepository struct {
	*BaseRepository[models.StationConfiguration]
	db  *gorm.DB
	log zerolog.Logger
}

func NewStationConfigRepository(db *gorm.DB) *StationConfigurationRepository {
	baseRepository := &BaseRepository[models.StationConfiguration]{
		DB:         db,
		Log:        logger.GetLogger("station-configuration-repository"),
		EntityName: "station-configuration-repository",
	}

	return &StationConfigurationRepository{
		BaseRepository: baseRepository,
		db:             db,
		log:            logger.GetLogger("station-config-repository"),
	}
}

func (s *StationConfigurationRepository) FindByStationId(ctx context.Context, stationId uint, includes map[string]bool) (*models.StationConfiguration, error) {
	var stationConfig models.StationConfiguration
	result := s.db.WithContext(ctx).Where("station_id = ?", stationId).First(&stationConfig)

	if result.Error != nil {
		return nil, result.Error
	}

	return &stationConfig, result.Error
}
