package services

import (
	"context"
	"github.com/rs/zerolog"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/core/interfaces"
	"gps-no-server/internal/infrastructure/http/dto"
)

type BaseService[T interfaces.Entity] struct {
	Repository interfaces.BaseRepository[T]
	Log        zerolog.Logger
	EntityName string
}

func NewBaseService[T interfaces.Entity](
	repository interfaces.BaseRepository[T],
	entityName string,
) *BaseService[T] {
	return &BaseService[T]{
		Repository: repository,
		Log:        logger.GetLogger(entityName + "-service"),
		EntityName: entityName,
	}
}

func (s *BaseService[T]) GetAll(ctx context.Context, includeParam *string) ([]*T, error) {
	includes := dto.ParseIncludes(includeParam)
	return s.Repository.FindAll(ctx, includes)
}

func (s *BaseService[T]) GetById(ctx context.Context, id uint, includeParam *string) (*T, error) {
	includes := dto.ParseIncludes(includeParam)
	return s.Repository.FindById(ctx, id, includes)
}

func (s *BaseService[T]) Create(ctx context.Context, entity *T, includeParam *string) (*T, error) {
	includes := dto.ParseIncludes(includeParam)
	return s.Repository.Create(ctx, entity, includes)
}

func (s *BaseService[T]) Update(ctx context.Context, entity *T, includeParam *string) (*T, error) {
	includes := dto.ParseIncludes(includeParam)
	return s.Repository.Update(ctx, entity, includes)
}

func (s *BaseService[T]) UpdateFields(ctx context.Context, entity *T, fields []string, includeParam *string) (*T, error) {
	includes := dto.ParseIncludes(includeParam)

	return s.Repository.UpdateFields(ctx, entity, fields, includes)
}

func (s *BaseService[T]) Delete(ctx context.Context, entity *T, includeParam *string) error {
	includes := dto.ParseIncludes(includeParam)
	return s.Repository.Delete(ctx, entity, includes)
}
