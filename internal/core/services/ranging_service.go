package services

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/core/repositories"
	"gps-no-server/internal/infrastructure/http/dto"
)

type RangingService struct {
	rangingRepository *repositories.RangingRepository
	stationService    *StationService
	eventPublisher    *RangingEventPublisher
	log               zerolog.Logger
}

func NewRangingService(rangingRepository *repositories.RangingRepository, stationService *StationService, eventStreamService *EventStreamService) *RangingService {
	service := &RangingService{
		rangingRepository: rangingRepository,
		stationService:    stationService,
		log:               logger.GetLogger("ranging-service"),
	}

	if eventStreamService != nil {
		service.eventPublisher = NewRangingEventPublisher(eventStreamService)
	}

	return service
}

func (s *RangingService) GetAll(ctx context.Context, sourceIdentifier string, destinationIdentifier string, includeParam *string) ([]*models.Ranging, error) {
	if sourceIdentifier != "" || destinationIdentifier != "" {
		return s.GetBySourceOrDestination(ctx, sourceIdentifier, destinationIdentifier, includeParam)
	}

	return s.rangingRepository.FindAll(ctx, true)
}

func (s *RangingService) GetById(ctx context.Context, id uint) (*models.Ranging, error) {
	return s.rangingRepository.FindById(ctx, id)
}

func (s *RangingService) GetByMac(ctx context.Context, mac string) ([]*models.Ranging, error) {
	return s.rangingRepository.FindByMac(ctx, mac)

}

func (s *RangingService) Save(ctx context.Context, ranging *models.Ranging) (*models.Ranging, error) {
	err := s.rangingRepository.Save(ctx, ranging)
	if err != nil {
		return nil, err
	}

	return ranging, nil
}

func (s *RangingService) SaveAll(ctx context.Context, rangingList []*models.Ranging) ([]*models.Ranging, error) {
	if len(rangingList) == 0 {
		return nil, nil
	}

	savedRangings, err := s.rangingRepository.SaveAll(ctx, rangingList)
	if err != nil {
		return nil, err
	}

	if s.eventPublisher != nil {
		for _, ranging := range savedRangings {
			if err := s.eventPublisher.PublishRangingEvent(ctx, ranging); err != nil {
				s.log.Error().Err(err).Msg("failed to publish ranging event")
			}
		}
	}

	return savedRangings, nil
}

func (s *RangingService) GetBySourceOrDestination(ctx context.Context, sourceIdentifier string, destinationIdentifier string, includeParam *string) ([]*models.Ranging, error) {
	var sourceStation *models.Station
	var destStation *models.Station
	var err error
	includes := dto.ParseIncludes(includeParam)

	if sourceIdentifier != "" {
		sourceStation, err = s.stationService.GetStationByIdentifier(ctx, sourceIdentifier, nil)
		if err != nil {
			return nil, fmt.Errorf("error while fetching source station: %w", err)
		}
	}

	if destinationIdentifier != "" {
		destStation, err = s.stationService.GetStationByIdentifier(ctx, destinationIdentifier, nil)

		if err != nil {
			return nil, fmt.Errorf("error while fetching destination station: %w", err)
		}
	}

	return s.rangingRepository.FindBySourceAndDestination(ctx, sourceStation, destStation, includes)
}
