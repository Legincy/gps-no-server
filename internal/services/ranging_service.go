package services

import (
	"context"
	"gps-no-server/internal/cache"
	"gps-no-server/internal/models"
	"gps-no-server/internal/repository"
)

type RangingService struct {
	rangingRepository *repository.RangingRepository
	rangingCache      *cache.RangingCache
	stationRepository *repository.StationRepository
}

func NewRangingService(rangingRepository *repository.RangingRepository, stationRepository *repository.StationRepository, cacheManager *cache.CacheManager) *RangingService {
	return &RangingService{
		rangingRepository: rangingRepository,
		rangingCache:      cacheManager.RangingCache,
		stationRepository: stationRepository,
	}
}

func (s *RangingService) GetAll(ctx context.Context, preloadTable bool) ([]*models.Ranging, error) {
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
