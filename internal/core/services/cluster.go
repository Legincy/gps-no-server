package services

import (
	"context"
	"github.com/rs/zerolog"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/core/repositories"
	"gps-no-server/internal/core/validation"
	"gps-no-server/internal/infrastructure/http/dto"
)

type ClusterService struct {
	*BaseService[models.Cluster]
	clusterRepository *repositories.ClusterRepository
	clusterValidator  *validation.ClusterValidator
	log               zerolog.Logger
}

func NewClusterService(clusterRepository *repositories.ClusterRepository) *ClusterService {
	baseService := NewBaseService[models.Cluster](
		clusterRepository,
		"cluster",
	)

	return &ClusterService{
		BaseService:       baseService,
		clusterRepository: clusterRepository,
		clusterValidator:  validation.NewClusterValidator(clusterRepository),
		log:               logger.GetLogger("cluster-service"),
	}
}

func (c *ClusterService) GetByMac(ctx context.Context, mac string, includeParam *string) (*models.Cluster, error) {
	includes := dto.ParseIncludes(includeParam)

	return c.clusterRepository.FindByMac(ctx, mac, includes)
}
