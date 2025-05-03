package services

import (
	"context"
	"gps-no-server/internal/models"
	"gps-no-server/internal/repository"
)

type ClusterService struct {
	clusterRepository *repository.ClusterRepository
}

func NewClusterService(clusterRepository *repository.ClusterRepository) *ClusterService {
	return &ClusterService{
		clusterRepository: clusterRepository,
	}
}

func (c *ClusterService) GetAll(ctx context.Context, preloadTable bool) ([]*models.Cluster, error) {
	return c.clusterRepository.FindAll(ctx, preloadTable)
}
