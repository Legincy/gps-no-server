package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/core/repositories"
	"gps-no-server/internal/infrastructure/http/dto"
)

type StationService struct {
	*BaseService[models.Station]
	stationRepository *repositories.StationRepository
	log               zerolog.Logger
}

func NewStationService(stationRepository *repositories.StationRepository) *StationService {
	baseService := NewBaseService[models.Station](
		stationRepository,
		"station",
	)

	return &StationService{
		BaseService:       baseService,
		stationRepository: stationRepository,
		log:               logger.GetLogger("services-station"),
	}
}

func (s *StationService) GetByMac(ctx context.Context, mac string, includeParam *string) (*models.Station, error) {
	includes := dto.ParseIncludes(includeParam)
	return s.stationRepository.FindByMac(ctx, mac, includes)
}

func (s *StationService) GetByIdentifier(ctx context.Context, identifier string, includeParam *string) (*models.Station, error) {
	includes := dto.ParseIncludes(includeParam)
	return s.stationRepository.FindByIdentifier(ctx, identifier, includes)
}

func (s *StationService) UpdateOrCreate(ctx context.Context, station *models.Station, includeParam *string) (*models.Station, error) {
	includes := dto.ParseIncludes(includeParam)

	existingStation, err := s.stationRepository.FindByMac(ctx, station.MacAddress, includes)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			createdStation, err := s.stationRepository.Create(ctx, station, includes)
			if err != nil {
				return nil, fmt.Errorf("failed to create station: %w", err)
			}
			return createdStation, nil
		}
		return nil, fmt.Errorf("error finding station: %w", err)
	}

	if existingStation != nil {

		return existingStation, nil
	}

	return nil, fmt.Errorf("unexpected condition in UpdateOrCreate")
}
