package validation

import (
	"context"
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/core/repositories"
)

type ClusterValidator struct {
	clusterRepository *repositories.ClusterRepository
}

func NewClusterValidator(clusterRepository *repositories.ClusterRepository) *ClusterValidator {
	return &ClusterValidator{
		clusterRepository: clusterRepository,
	}
}

func (c *ClusterValidator) ValidateCreate(ctx context.Context, cluster *models.Cluster) error {
	return c.validateCluster(ctx, cluster, true)
}

func (c *ClusterValidator) validateCluster(ctx context.Context, cluster *models.Cluster, isNew bool) error {
	var errors ValidationErrors

	if cluster.Name == "" {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "Name cannot be empty",
		})
	}

	if cluster.Name != "" && isNew {
		existingClusters, err := c.clusterRepository.FindAll(ctx, map[string]bool{})
		if err == nil {
			for _, existingCluster := range existingClusters {
				if existingCluster.Name == cluster.Name && (isNew || existingCluster.ID != cluster.ID) {
					errors = append(errors, ValidationError{
						Field:   "name",
						Message: "A cluster with this name already exists",
					})
					break
				}
			}
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
