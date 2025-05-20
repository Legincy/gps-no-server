package services

import (
	"context"
	"github.com/rs/zerolog"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/core/repositories"
	"gps-no-server/internal/infrastructure/http/dto"
)

type StationConfigurationService struct {
	*BaseService[models.StationConfiguration]
	stationConfigurationRepository *repositories.StationConfigurationRepository
	log                            zerolog.Logger
}

func NewStationConfigService(stationConfigRepository *repositories.StationConfigurationRepository) *StationConfigurationService {
	baseService := NewBaseService[models.StationConfiguration](
		stationConfigRepository,
		"station-configuration",
	)

	return &StationConfigurationService{
		BaseService:                    baseService,
		stationConfigurationRepository: stationConfigRepository,
		log:                            logger.GetLogger("station-configuration-service"),
	}
}

func (s *StationConfigurationService) GetByStationId(ctx context.Context, stationId uint, includeParam *string) (*models.StationConfiguration, error) {
	includes := dto.ParseIncludes(includeParam)
	return s.stationConfigurationRepository.FindByStationId(ctx, stationId, includes)
}
