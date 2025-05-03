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

func (s *StationService) GetAllStations(ctx context.Context, preloadTable bool) ([]*models.Station, error) {
	return s.stationRepository.FindAll(ctx, preloadTable)
}

func (s *StationService) GetStationByID(ctx context.Context, id uint) (*models.Station, error) {
	return s.stationRepository.FindByID(ctx, id)
}

func (s *StationService) GetStationByMac(ctx context.Context, mac string) (*models.Station, error) {
	return s.stationRepository.FindByMac(ctx, mac)
}

func (s *StationService) GetActiveStations(ctx context.Context) ([]*models.Station, error) {
	return s.stationRepository.FindActive(ctx)
}

func (s *StationService) SaveStation(ctx context.Context, station *models.Station) (*models.Station, error) {
	return s.stationRepository.Save(ctx, station)
}
