package interfaces

import (
	"context"
	"gps-no-server/internal/models"
)

type StationService interface {
	GetAllStations(ctx context.Context) ([]*models.Station, error)
	GetStationByID(ctx context.Context, id uint) (*models.Station, error)
	GetStationByMac(ctx context.Context, mac string) (*models.Station, error)
	GetActiveStations(ctx context.Context) ([]*models.Station, error)
	SaveStation(ctx context.Context, station *models.Station) (*models.Station, error)
}

type RangingService interface {
	GetAllRanging(ctx context.Context) ([]*models.Ranging, error)
	GetRangingByID(ctx context.Context, id uint) (*models.Ranging, error)
	GetRangingByMac(ctx context.Context, mac string) ([]*models.Ranging, error)
	GetRangingByStationID(ctx context.Context, stationID uint) ([]*models.Ranging, error)
	SaveRanging(ctx context.Context, ranging *models.Ranging) (*models.Ranging, error)
}
