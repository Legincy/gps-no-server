package services

import (
	"context"
	"gps-no-server/internal/models"
	"gps-no-server/internal/repository"
)

type stationService struct {
	stationRepository *repository.StationRepository
}

func NewStationService(stationRepository *repository.StationRepository) *stationService {
	return &stationService{
		stationRepository: stationRepository,
	}
}

func (s *stationService) GetAllStations(ctx context.Context) ([]*models.Station, error) {
	return s.stationRepository.FindAll(ctx)
}

func (s *stationService) GetStationByID(ctx context.Context, id uint) (*models.Station, error) {
	return s.stationRepository.FindByID(ctx, id)
}

func (s *stationService) GetStationByMac(ctx context.Context, mac string) (*models.Station, error) {
	return s.stationRepository.FindByMac(ctx, mac)
}

func (s *stationService) GetActiveStations(ctx context.Context) ([]*models.Station, error) {
	return s.stationRepository.FindActive(ctx)
}

func (s *stationService) SaveStation(ctx context.Context, station *models.Station) (*models.Station, error) {
	return s.stationRepository.Save(ctx, station)
}
