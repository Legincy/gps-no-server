package services

import (
	"context"
	"github.com/rs/zerolog"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/di/interfaces"
)

type ClusterService struct {
	clusterRepository interfaces.ClusterRepository
	log               zerolog.Logger
}

func NewClusterService(clusterRepository interfaces.ClusterRepository) *ClusterService {
	return &ClusterService{
		clusterRepository: clusterRepository,
		log:               logger.GetLogger("cluster-service"),
	}
}

func (c *ClusterService) GetAll(ctx context.Context) ([]*models.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ClusterService) GetByID(ctx context.Context, id uint) (*models.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ClusterService) Create(ctx context.Context, cluster *models.Cluster) (*models.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ClusterService) Update(ctx context.Context, cluster *models.Cluster) (*models.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ClusterService) Delete(ctx context.Context, id uint) error {
	//TODO implement me
	panic("implement me")
}
