package interfaces

import "context"

type BaseRepository[T any] interface {
	FindAll(ctx context.Context, includes map[string]bool) ([]*T, error)
	FindById(ctx context.Context, id uint, includes map[string]bool) (*T, error)
	Create(ctx context.Context, entity *T, includes map[string]bool) (*T, error)
	Update(ctx context.Context, entity *T, includes map[string]bool) (*T, error)
	UpdateFields(ctx context.Context, entity *T, fields []string, includes map[string]bool) (*T, error)
	Delete(ctx context.Context, entity *T, includes map[string]bool) error
}
