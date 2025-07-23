package services

import (
	"context"
	"github.com/rs/zerolog"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/di/interfaces"
)

type StationService struct {
	stationRepository interfaces.StationRepository
	log               zerolog.Logger
}

func NewStationService(stationRepository interfaces.StationRepository) *StationService {
	return &StationService{
		stationRepository: stationRepository,
		log:               logger.GetLogger("services-station"),
	}
}

func (s *StationService) GetAll(ctx context.Context) ([]*models.Station, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StationService) GetByID(ctx context.Context, id uint) (*models.Station, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StationService) GetByMac(ctx context.Context, mac string) (*models.Station, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StationService) Create(ctx context.Context, station *models.Station) (*models.Station, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StationService) Update(ctx context.Context, station *models.Station) (*models.Station, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StationService) Delete(ctx context.Context, id uint) error {
	//TODO implement me
	panic("implement me")
}
