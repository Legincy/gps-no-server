package services

import (
	"context"
	"gps-no-server/internal/models"
	"gps-no-server/internal/repository"
)

type StationService struct {
	stationRepository *repository.StationRepository
}

func NewStationService(stationRepository *repository.StationRepository) *StationService {
	return &StationService{
		stationRepository: stationRepository,
	}
}

func (s *StationService) GetAll(ctx context.Context, preloadTable bool) ([]*models.Station, error) {
	return s.stationRepository.FindAll(ctx, preloadTable)
}

func (s *StationService) GetByID(ctx context.Context, id uint) (*models.Station, error) {
	return s.stationRepository.FindByID(ctx, id)
}

func (s *StationService) GetByMac(ctx context.Context, mac string) (*models.Station, error) {
	return s.stationRepository.FindByMac(ctx, mac)
}

func (s *StationService) GetActive(ctx context.Context) ([]*models.Station, error) {
	return s.stationRepository.FindActive(ctx)
}

func (s *StationService) Save(ctx context.Context, station *models.Station) (*models.Station, error) {
	return s.stationRepository.Save(ctx, station)
}
