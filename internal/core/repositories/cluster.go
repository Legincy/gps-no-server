package repositories

import (
	"context"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/core/models"
)

type ClusterRepository struct {
	*BaseRepository[models.Cluster]
	db  *gorm.DB
	log zerolog.Logger
}

func NewClusterRepository(db *gorm.DB) *ClusterRepository {
	baseRepository := &BaseRepository[models.Cluster]{
		DB:         db,
		Log:        logger.GetLogger("cluster-repository"),
		EntityName: "cluster-repository",
	}

	return &ClusterRepository{
		BaseRepository: baseRepository,
		db:             db,
		log:            logger.GetLogger("cluster-repository"),
	}
}

func (c *ClusterRepository) FindByMac(ctx context.Context, macAddress string, includes map[string]bool) (*models.Cluster, error) {
	var cluster models.Cluster
	result := c.db.WithContext(ctx).Where("mac_address = ?", macAddress).First(&cluster)

	if result.Error != nil {
		return nil, result.Error
	}

	return &cluster, result.Error
}

func (c *ClusterRepository) FindAll(ctx context.Context) ([]*models.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ClusterRepository) FindByID(ctx context.Context, id uint) (*models.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ClusterRepository) Create(ctx context.Context, cluster *models.Cluster) (*models.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ClusterRepository) Update(ctx context.Context, cluster *models.Cluster) (*models.Cluster, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ClusterRepository) Delete(ctx context.Context, id uint) error {
	//TODO implement me
	panic("implement me")
}
