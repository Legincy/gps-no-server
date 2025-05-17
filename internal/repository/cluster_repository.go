package repository

import (
	"context"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gps-no-server/internal/logger"
	"gps-no-server/internal/models"
)

type ClusterRepository struct {
	db  *gorm.DB
	log zerolog.Logger
}

func NewClusterRepository(db *gorm.DB) *ClusterRepository {
	return &ClusterRepository{
		db:  db,
		log: logger.GetLogger("cluster-repository"),
	}
}

func (c *ClusterRepository) Save(ctx context.Context, cluster *models.Cluster, preloadTable bool) (*models.Cluster, error) {
	result := c.db.WithContext(ctx).Where("id = ?", cluster.ID).FirstOrCreate(cluster)

	return cluster, result.Error
}

func (c *ClusterRepository) FindAll(ctx context.Context, preloadTable bool) ([]*models.Cluster, error) {
	var clusters []*models.Cluster

	query := c.db.WithContext(ctx)

	if preloadTable {
		query = query.Preload("Stations")
	}

	result := query.Find(&clusters)

	return clusters, result.Error
}

func (c *ClusterRepository) FindById(ctx context.Context, id uint) (*models.Cluster, error) {
	var cluster models.Cluster
	result := c.db.WithContext(ctx).Where("id = ?", id).First(&cluster)

	if result.Error != nil {
		return nil, result.Error
	}

	return &cluster, result.Error
}

func (c *ClusterRepository) FindByMac(ctx context.Context, macAddress string) (*models.Cluster, error) {
	var cluster models.Cluster
	result := c.db.WithContext(ctx).Where("mac_address = ?", macAddress).First(&cluster)

	if result.Error != nil {
		return nil, result.Error
	}

	return &cluster, result.Error
}

func (c *ClusterRepository) Create(ctx context.Context, cluster *models.Cluster) (*models.Cluster, error) {
	result := c.db.WithContext(ctx).Create(cluster)

	if result.Error != nil {
		return nil, result.Error
	}

	return cluster, nil
}

func (c *ClusterRepository) Update(ctx context.Context, cluster *models.Cluster) (*models.Cluster, error) {
	result := c.db.WithContext(ctx).Save(cluster)

	if result.Error != nil {
		return nil, result.Error
	}

	return cluster, nil
}

func (c *ClusterRepository) Delete(ctx context.Context, cluster *models.Cluster) error {
	result := c.db.WithContext(ctx).Delete(cluster)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (c *ClusterRepository) DeleteById(ctx context.Context, id uint) error {
	var cluster models.Cluster
	result := c.db.WithContext(ctx).Where("id = ?", id).Delete(&cluster)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
