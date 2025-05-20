package repositories

import (
	"context"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/core/interfaces"
	"reflect"
)

type BaseRepository[T interfaces.Entity] struct {
	DB         *gorm.DB
	Log        zerolog.Logger
	EntityName string
}

func NewBaseRepository[T interfaces.Entity](
	db *gorm.DB,
	entityName string,
) *BaseRepository[T] {
	return &BaseRepository[T]{
		DB:         db,
		Log:        logger.GetLogger(entityName + "-repository"),
		EntityName: entityName,
	}
}

func (r *BaseRepository[T]) FindAll(ctx context.Context, includes map[string]bool) ([]*T, error) {
	var entities []*T
	query := r.DB.WithContext(ctx)

	result := query.Find(&entities)
	return entities, result.Error
}

func (r *BaseRepository[T]) FindById(ctx context.Context, id uint, includes map[string]bool) (*T, error) {
	var entity *T
	query := r.DB.WithContext(ctx)

	result := query.First(&entity, id)
	return entity, result.Error
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T, includes map[string]bool) (*T, error) {
	result := r.DB.WithContext(ctx).Create(&entity)
	return entity, result.Error
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T, includes map[string]bool) (*T, error) {
	tx := r.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return entity, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(entity).Error; err != nil {
		tx.Rollback()
		return entity, err
	}

	id := r.getID(entity)

	updatedEntity := entity

	if err := tx.First(&updatedEntity, id).Error; err != nil {
		tx.Rollback()
		return entity, err
	}

	if len(includes) > 0 {
		query := tx

		if err := query.First(&updatedEntity, id).Error; err != nil {
			tx.Rollback()
			return entity, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return entity, err
	}

	return updatedEntity, nil
}

func (r *BaseRepository[T]) UpdateFields(ctx context.Context, entity *T, fields []string, includes map[string]bool) (*T, error) {
	if len(fields) == 0 {
		return r.Update(ctx, entity, includes)
	}

	tx := r.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return entity, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(entity).Select(fields).Updates(entity).Error; err != nil {
		tx.Rollback()
		return entity, err
	}

	id := r.getID(entity)
	updatedEntity := entity

	if err := tx.First(&updatedEntity, id).Error; err != nil {
		tx.Rollback()
		return entity, err
	}

	if len(includes) > 0 {
		query := tx

		if err := query.First(&updatedEntity, id).Error; err != nil {
			tx.Rollback()
			return entity, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return entity, err
	}

	return updatedEntity, nil
}

func (r *BaseRepository[T]) getID(entity *T) uint {
	val := reflect.ValueOf(entity)

	idField := val.FieldByName("ID")
	if !idField.IsValid() {
		return 0
	}

	return uint(idField.Uint())
}

func (r *BaseRepository[T]) Delete(ctx context.Context, entity *T, includes map[string]bool) error {
	result := r.DB.WithContext(ctx).Delete(entity)
	return result.Error
}

func (r *BaseRepository[T]) Save(ctx context.Context, entity *T, includes map[string]bool) (*T, error) {
	result := r.DB.WithContext(ctx).Save(entity)
	return entity, result.Error
}

func (r *BaseRepository[T]) UpdateOrCreate(ctx context.Context, entity *T, includes map[string]bool) (*T, error) {
	result := r.DB.WithContext(ctx).Save(entity)
	return entity, result.Error
}
