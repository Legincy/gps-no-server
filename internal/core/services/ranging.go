package services

import (
	"context"
	"github.com/rs/zerolog"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/core/repositories"
	"gps-no-server/internal/infrastructure/http/dto"
)

type RangingService struct {
	*BaseService[models.Ranging]
	rangingRepository *repositories.RangingRepository
	stationService    *StationService
	eventPublisher    *RangingEventPublisher
	log               zerolog.Logger
}

func NewRangingService(rangingRepository *repositories.RangingRepository, stationService *StationService, eventStreamService *EventStreamService) *RangingService {
	baseService := NewBaseService[models.Ranging](
		rangingRepository,
		"ranging",
	)

	service := &RangingService{
		BaseService:       baseService,
		rangingRepository: rangingRepository,
		stationService:    stationService,
		log:               logger.GetLogger("ranging-service"),
	}

	if eventStreamService != nil {
		service.eventPublisher = NewRangingEventPublisher(eventStreamService)
	}

	return service
}

func (s *RangingService) GetByMac(ctx context.Context, mac string, includeParam *string) ([]*models.Ranging, error) {
	includes := dto.ParseIncludes(includeParam)

	return s.rangingRepository.FindByMac(ctx, mac, includes)
}

func (s *RangingService) GetBySourceStationAndDestinationStation(ctx context.Context, source *models.Station, destination *models.Station, includes map[string]bool) (*models.Ranging, error) {
	ranging, err := s.rangingRepository.FindBySourceStationAndDestinationStation(ctx, source, destination, includes)
	if err != nil {
		return nil, err
	}

	return ranging, nil
}
