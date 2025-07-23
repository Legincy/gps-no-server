package interfaces

import (
	"context"
	"gps-no-server/internal/core/models"
)

type StationRepository interface {
	FindAll(ctx context.Context) ([]*models.Station, error)
	FindByID(ctx context.Context, id uint) (*models.Station, error)
	FindByMac(ctx context.Context, mac string) (*models.Station, error)
	Create(ctx context.Context, station *models.Station) (*models.Station, error)
	Update(ctx context.Context, station *models.Station) (*models.Station, error)
	Delete(ctx context.Context, id uint) error
}

type ClusterRepository interface {
	FindAll(ctx context.Context) ([]*models.Cluster, error)
	FindByID(ctx context.Context, id uint) (*models.Cluster, error)
	Create(ctx context.Context, cluster *models.Cluster) (*models.Cluster, error)
	Update(ctx context.Context, cluster *models.Cluster) (*models.Cluster, error)
	Delete(ctx context.Context, id uint) error
}

type RangingRepository interface {
	FindAll(ctx context.Context) ([]*models.Ranging, error)
	FindByID(ctx context.Context, id uint) (*models.Ranging, error)
	Create(ctx context.Context, ranging *models.Ranging) (*models.Ranging, error)
	Update(ctx context.Context, ranging *models.Ranging) (*models.Ranging, error)
	Delete(ctx context.Context, id uint) error
}
