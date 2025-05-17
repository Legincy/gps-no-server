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

func (c *ClusterService) GetById(ctx context.Context, id uint) (*models.Cluster, error) {
	return c.clusterRepository.FindById(ctx, id)
}

func (c *ClusterService) GetByMac(ctx context.Context, mac string) (*models.Cluster, error) {
	return c.clusterRepository.FindByMac(ctx, mac)
}

func (c *ClusterService) Save(ctx context.Context, cluster *models.Cluster) (*models.Cluster, error) {
	return c.clusterRepository.Save(ctx, cluster, false)
}

func (c *ClusterService) SaveAll(ctx context.Context, clusterList []*models.Cluster) ([]*models.Cluster, error) {
	if len(clusterList) == 0 {
		return nil, nil
	}

	for _, cluster := range clusterList {
		_, err := c.clusterRepository.Save(ctx, cluster, false)
		if err != nil {
			return nil, err
		}
	}

	return clusterList, nil
}

func (c *ClusterService) Create(ctx context.Context, cluster *models.Cluster) (*models.Cluster, error) {
	result, err := c.clusterRepository.Create(ctx, cluster)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *ClusterService) Delete(ctx context.Context, cluster *models.Cluster) error {
	return c.clusterRepository.DeleteById(ctx, cluster.ID)
}

func (c *ClusterService) Update(ctx context.Context, cluster *models.Cluster) (*models.Cluster, error) {
	result, err := c.clusterRepository.Update(ctx, cluster)
	if err != nil {
		return nil, err
	}

	return result, nil
}
