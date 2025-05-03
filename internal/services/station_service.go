package services

import (
	"context"
	"fmt"
	"gps-no-server/internal/models"
	"gps-no-server/internal/repository"
	"strconv"
	"strings"
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

func (s *StationService) GetStationByIdentifier(ctx context.Context, identifier string) (*models.Station, error) {
	if strings.Contains(identifier, ":") {
		return s.stationRepository.FindByMac(ctx, identifier)
	}

	id, err := strconv.ParseUint(identifier, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("ungültiger Identifier '%s': weder MAC-Adresse noch gültige ID", identifier)
	}

	return s.stationRepository.FindByID(ctx, uint(id))
}
