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

func (c *RangingRepository) Save(ctx context.Context, ranging *models.Ranging) error {
	result := c.db.WithContext(ctx).Where(" source_id = ? AND target_id = ?", ranging.Source.ID, ranging.Destination.ID).Assign(ranging).FirstOrCreate(ranging)

	return result.Error
}

func (c *RangingRepository) SaveAll(ctx context.Context, rangingList []*models.Ranging) error {
	if len(rangingList) == 0 {
		return nil
	}

	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, ranging := range rangingList {
		var existingRanging models.Ranging
		result := tx.Where("source_id = ? AND destination_id = ? AND ABS(EXTRACT(EPOCH FROM (timestamp - ?))) < 1",
			ranging.Source.ID, ranging.Destination.ID, ranging.Timestamp).
			First(&existingRanging)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				if err := tx.Create(ranging).Error; err != nil {
					tx.Rollback()
					return err
				}
			} else {
				tx.Rollback()
				return result.Error
			}
		} else {
			existingRanging.RawDistance = ranging.RawDistance

			if err := tx.Save(&existingRanging).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}
