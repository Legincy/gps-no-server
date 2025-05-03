package services

import (
	"context"
	"fmt"
	"gps-no-server/internal/cache"
	"gps-no-server/internal/models"
	"gps-no-server/internal/repository"
)

type RangingService struct {
	rangingRepository *repository.RangingRepository
	rangingCache      *cache.RangingCache
	stationService    *StationService
}

func NewRangingService(rangingRepository *repository.RangingRepository, stationService *StationService, cacheManager *cache.CacheManager) *RangingService {
	return &RangingService{
		rangingRepository: rangingRepository,
		rangingCache:      cacheManager.RangingCache,
		stationService:    stationService,
	}
}

func (s *RangingService) GetAll(ctx context.Context, preloadTable bool, sourceIdentifier string, destinationIdentifier string) ([]*models.Ranging, error) {
	if sourceIdentifier != "" || destinationIdentifier != "" {
		return s.GetBySourceOrDestination(ctx, preloadTable, sourceIdentifier, destinationIdentifier)
	}

	return s.rangingRepository.FindAll(ctx, preloadTable)
}

func (s *RangingService) GetByID(ctx context.Context, id uint) (*models.Ranging, error) {
	return s.rangingRepository.FindByID(ctx, id)
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

func (s *RangingService) SaveAll(ctx context.Context, rangingList []*models.Ranging) error {
	if len(rangingList) == 0 {
		return nil
	}

	err := s.rangingRepository.SaveAll(ctx, rangingList)
	if err != nil {
		return err
	}

	return nil
}

func (s *RangingService) GetBySourceOrDestination(ctx context.Context, preloadTable bool, sourceIdentifier string, destinationIdentifier string) ([]*models.Ranging, error) {
	var sourceStation *models.Station
	var destStation *models.Station
	var err error

	if sourceIdentifier != "" {
		sourceStation, err = s.stationService.GetStationByIdentifier(ctx, sourceIdentifier)
		if err != nil {
			return nil, fmt.Errorf("error while fetching source station: %w", err)
		}
	}

	if destinationIdentifier != "" {
		destStation, err = s.stationService.GetStationByIdentifier(ctx, destinationIdentifier)

		if err != nil {
			return nil, fmt.Errorf("error while fetching destination station: %w", err)
		}
	}

	return s.rangingRepository.FindBySourceAndDestination(ctx, preloadTable, sourceStation, destStation)
}
