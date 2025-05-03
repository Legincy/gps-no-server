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
