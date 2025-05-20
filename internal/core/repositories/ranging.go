package repositories

import (
	"context"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/core/models"
)

type RangingRepository struct {
	*BaseRepository[models.Ranging]
	db  *gorm.DB
	log zerolog.Logger
}

func NewRangingRepository(db *gorm.DB) *RangingRepository {
	baseRepository := &BaseRepository[models.Ranging]{
		DB:         db,
		Log:        logger.GetLogger("ranging-repository"),
		EntityName: "ranging-repository",
	}

	return &RangingRepository{
		BaseRepository: baseRepository,
		db:             db,
		log:            logger.GetLogger("ranging-repository"),
	}
}

func (r *RangingRepository) FindByMac(ctx context.Context, mac string, includes map[string]bool) ([]*models.Ranging, error) {
	var rangings []*models.Ranging
	query := r.db.WithContext(ctx).Where("mac = ?", mac)

	if includes["station"] {
		query = query.Preload("Station")
	}

	result := query.Find(&rangings)

	return rangings, result.Error
}

func (r *RangingRepository) FindBySourceStationAndDestinationStation(ctx context.Context, source *models.Station, destination *models.Station, includes map[string]bool) (*models.Ranging, error) {
	var ranging models.Ranging
	result := r.db.WithContext(ctx).Where("source_id = ? AND destination_id = ?", source.ID, destination.ID).First(&ranging)

	if result.Error != nil {
		return nil, result.Error
	}

	return &ranging, result.Error
}
