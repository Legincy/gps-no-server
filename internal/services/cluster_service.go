package services

import (
	"context"
	"github.com/rs/zerolog"
	"gps-no-server/internal/controllers/dto"
	"gps-no-server/internal/logger"
	"gps-no-server/internal/models"
	"gps-no-server/internal/repository"
)

type ClusterService struct {
	clusterRepository *repository.ClusterRepository
	log               zerolog.Logger
}

func NewClusterService(clusterRepository *repository.ClusterRepository) *ClusterService {
	return &ClusterService{
		clusterRepository: clusterRepository,
		log:               logger.GetLogger("cluster-service"),
	}
}

func (c *ClusterService) GetAll(ctx context.Context, includeParam *string) ([]*models.Cluster, error) {
	includes := dto.ParseIncludes(includeParam)

	return c.clusterRepository.FindAll(ctx, includes)
}

func (c *ClusterService) GetById(ctx context.Context, id uint, includeParam *string) (*models.Cluster, error) {
	includes := dto.ParseIncludes(includeParam)

	return c.clusterRepository.FindById(ctx, id, includes)
}

func (c *ClusterService) GetByMac(ctx context.Context, mac string, includeParam *string) (*models.Cluster, error) {
	includes := dto.ParseIncludes(includeParam)

	return c.clusterRepository.FindByMac(ctx, mac, includes)
}

func (c *ClusterService) Save(ctx context.Context, cluster *models.Cluster, includeParam *string) (*models.Cluster, error) {
	includes := dto.ParseIncludes(includeParam)

	return c.clusterRepository.Save(ctx, cluster, includes)
}

func (c *ClusterService) SaveAll(ctx context.Context, clusterList []*models.Cluster, includeParam *string) ([]*models.Cluster, error) {
	if len(clusterList) == 0 {
		return nil, nil
	}
	includes := dto.ParseIncludes(includeParam)

	for _, cluster := range clusterList {
		_, err := c.clusterRepository.Save(ctx, cluster, includes)
		if err != nil {
			return nil, err
		}
	}

	return clusterList, nil
}

func (c *ClusterService) Create(ctx context.Context, cluster *models.Cluster, includeParam *string) (*models.Cluster, error) {
	includes := dto.ParseIncludes(includeParam)

	result, err := c.clusterRepository.Create(ctx, cluster, includes)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *ClusterService) Delete(ctx context.Context, cluster *models.Cluster, includeParam *string) error {
	includes := dto.ParseIncludes(includeParam)

	return c.clusterRepository.DeleteById(ctx, cluster.ID, includes)
}

func (c *ClusterService) Update(ctx context.Context, cluster *models.Cluster, includeParam *string) (*models.Cluster, error) {
	includes := dto.ParseIncludes(includeParam)

	result, err := c.clusterRepository.Update(ctx, cluster, includes)
	if err != nil {
		return nil, err
	}

	return result, nil
}
