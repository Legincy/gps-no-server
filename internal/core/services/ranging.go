package services

import (
	"context"
	"github.com/rs/zerolog"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/di/interfaces"
)

type RangingService struct {
	rangingRepository interfaces.RangingRepository
	influx            interfaces.InfluxConnection
	log               zerolog.Logger
}

func NewRangingService(rangingRepository interfaces.RangingRepository, influx interfaces.InfluxConnection) *RangingService {
	service := &RangingService{
		rangingRepository: rangingRepository,
		influx:            influx,
		log:               logger.GetLogger("ranging-service"),
	}

	return service
}

func (s *RangingService) GetAll(ctx context.Context) ([]*models.Ranging, error) {
	//TODO implement me
	panic("implement me")
}

func (s *RangingService) GetByID(ctx context.Context, id uint) (*models.Ranging, error) {
	//TODO implement me
	panic("implement me")
}

func (s *RangingService) Create(ctx context.Context, ranging *models.Ranging) (*models.Ranging, error) {
	//TODO implement me
	panic("implement me")
}

func (s *RangingService) StoreMeasurement(ctx context.Context, ranging *models.Ranging) error {
	//TODO implement me
	panic("implement me")
}
