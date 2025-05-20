package interfaces

import (
	"context"
)

type CrudService[T Entity] interface {
	GetAll(ctx context.Context, includeParam *string) ([]T, error)
	GetById(ctx context.Context, id uint, includeParam *string) (T, error)
	Create(ctx context.Context, entity T, includeParam *string) (T, error)
	Update(ctx context.Context, entity T, includeParam *string) (T, error)
	UpdateFields(ctx context.Context, entity T, fields []string, includeParam *string) (T, error)
	Delete(ctx context.Context, entity T, includeParam *string) error
}
