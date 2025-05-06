package repository

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gps-no-server/internal/logger"
	"gps-no-server/internal/models"
)

type RangingRepository struct {
	db  *gorm.DB
	log zerolog.Logger
}

func NewRangingRepository(db *gorm.DB) *RangingRepository {
	return &RangingRepository{
		db:  db,
		log: logger.GetLogger("ranging-repository"),
	}
}

func (c *RangingRepository) FindAll(ctx context.Context, preloadTable bool) ([]*models.Ranging, error) {
	var rangings []*models.Ranging

	query := c.db.WithContext(ctx)

	if preloadTable {
		query = query.Preload("Source").Preload("Destination")
	}

	result := query.Find(&rangings)

	return rangings, result.Error
}

func (r *RangingRepository) FindByID(ctx context.Context, id uint) (*models.Ranging, error) {
	var ranging models.Ranging
	result := r.db.WithContext(ctx).First(&ranging, id)

	return &ranging, result.Error
}

func (r *RangingRepository) FindByMac(ctx context.Context, mac string) ([]*models.Ranging, error) {
	var rangings []*models.Ranging
	result := r.db.WithContext(ctx).Where("source_id = ? OR destination_id = ?", mac, mac).Find(&rangings)

	return rangings, result.Error
}

func (c *RangingRepository) Save(ctx context.Context, ranging *models.Ranging) error {
	result := c.db.WithContext(ctx).Where(" source_id = ? AND target_id = ?", ranging.Source.ID, ranging.Destination.ID).Assign(ranging).FirstOrCreate(ranging)

	return result.Error
}

func (c *RangingRepository) SaveAll(ctx context.Context, rangingList []*models.Ranging) ([]*models.Ranging, error) {
	if len(rangingList) == 0 {
		return nil, nil
	}

	query := c.db.WithContext(ctx).Begin()
	if query.Error != nil {
		return nil, query.Error
	}

	defer func() {
		if r := recover(); r != nil {
			query.Rollback()
		}
	}()

	savedRangings := make([]*models.Ranging, 0, len(rangingList))

	for _, ranging := range rangingList {
		var existingRanging models.Ranging
		result := query.Where("source_id = ? AND destination_id = ?",
			ranging.Source.ID, ranging.Destination.ID).
			First(&existingRanging)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				if err := query.Create(ranging).Error; err != nil {
					query.Rollback()
					return nil, err
				}
				savedRangings = append(savedRangings, ranging)
			} else {
				query.Rollback()
				return nil, result.Error
			}
		} else {
			existingRanging.RawDistance = ranging.RawDistance

			if err := query.Save(&existingRanging).Error; err != nil {
				query.Rollback()
				return nil, err
			}

			savedRangings = append(savedRangings, &existingRanging)
		}
	}

	if err := query.Commit().Error; err != nil {
		return nil, err
	}

	return savedRangings, nil
}

func (c *RangingRepository) FindBySourceAndDestination(ctx context.Context, preloadTable bool, sourceStation *models.Station, destinationStation *models.Station) ([]*models.Ranging, error) {
	var rangings []*models.Ranging
	query := c.db.WithContext(ctx)

	if preloadTable {
		query = query.Preload("Source").Preload("Destination")
	}

	if sourceStation != nil {
		query = query.Where("source_id = ?", sourceStation.ID)
	}

	if destinationStation != nil {
		query = query.Where("destination_id = ?", destinationStation.ID)
	}

	result := query.Find(&rangings)

	return rangings, result.Error
}
