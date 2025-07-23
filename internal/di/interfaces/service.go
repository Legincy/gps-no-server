package interfaces

import (
	"context"
	"gps-no-server/internal/core/models"
)

type StationService interface {
	GetAll(ctx context.Context) ([]*models.Station, error)
	GetByID(ctx context.Context, id uint) (*models.Station, error)
	GetByMac(ctx context.Context, mac string) (*models.Station, error)
	Create(ctx context.Context, station *models.Station) (*models.Station, error)
	Update(ctx context.Context, station *models.Station) (*models.Station, error)
	Delete(ctx context.Context, id uint) error
}

type ClusterService interface {
	GetAll(ctx context.Context) ([]*models.Cluster, error)
	GetByID(ctx context.Context, id uint) (*models.Cluster, error)
	Create(ctx context.Context, cluster *models.Cluster) (*models.Cluster, error)
	Update(ctx context.Context, cluster *models.Cluster) (*models.Cluster, error)
	Delete(ctx context.Context, id uint) error
}

type RangingService interface {
	GetAll(ctx context.Context) ([]*models.Ranging, error)
	GetByID(ctx context.Context, id uint) (*models.Ranging, error)
	Create(ctx context.Context, ranging *models.Ranging) (*models.Ranging, error)
	StoreMeasurement(ctx context.Context, ranging *models.Ranging) error
}
