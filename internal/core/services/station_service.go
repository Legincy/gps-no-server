package services

import (
	"context"
	"fmt"
	"gps-no-server/internal/core/models"
	"gps-no-server/internal/core/repositories"
	"gps-no-server/internal/events"
	"gps-no-server/internal/infrastructure/http/dto"
	"strconv"
	"strings"
	"time"
)

type StationService struct {
	stationRepository *repositories.StationRepository
	stationEventBus   *events.StationEventBus
}

func NewStationService(stationRepository *repositories.StationRepository, stationEventBus *events.StationEventBus) *StationService {
	return &StationService{
		stationRepository: stationRepository,
		stationEventBus:   stationEventBus,
	}
}

func (s *StationService) GetAll(ctx context.Context, includeParam *string) ([]*models.Station, error) {
	includes := dto.ParseIncludes(includeParam)

	return s.stationRepository.FindAll(ctx, includes)
}

func (s *StationService) GetById(ctx context.Context, id uint, includeParam *string) (*models.Station, error) {
	includes := dto.ParseIncludes(includeParam)

	return s.stationRepository.FindByID(ctx, id, includes)
}

func (s *StationService) GetByMac(ctx context.Context, mac string, includeParam *string) (*models.Station, error) {
	includes := dto.ParseIncludes(includeParam)

	return s.stationRepository.FindByMac(ctx, mac, includes)
}

func (s *StationService) GetActive(ctx context.Context, includeParam *string) ([]*models.Station, error) {
	return s.stationRepository.FindActive(ctx)
}

func (s *StationService) Save(ctx context.Context, station *models.Station, includeParam *string) (*models.Station, error) {
	return s.stationRepository.Save(ctx, station)
}

func (s *StationService) GetStationByIdentifier(ctx context.Context, identifier string, includeParam *string) (*models.Station, error) {
	includes := dto.ParseIncludes(includeParam)
	if strings.Contains(identifier, ":") {
		return s.stationRepository.FindByMac(ctx, identifier, includes)
	}

	id, err := strconv.ParseUint(identifier, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("ungültiger Identifier '%s': weder MAC-Adresse noch gültige ID", identifier)
	}

	return s.stationRepository.FindByID(ctx, uint(id), includes)
}

func (s *StationService) SaveAll(ctx context.Context, stationList []*models.Station, includeParam *string) ([]*models.Station, error) {
	if len(stationList) == 0 {
		return nil, nil
	}

	for _, station := range stationList {
		_, err := s.stationRepository.Save(ctx, station)
		if err != nil {
			return nil, err
		}
	}

	return stationList, nil
}

func (s *StationService) Delete(ctx context.Context, station *models.Station, includeParam *string) error {
	stationId := station.ID

	return s.stationRepository.DeleteById(ctx, stationId)
}

func (s *StationService) Update(ctx context.Context, station *models.Station, includeParam *string) (*models.Station, error) {
	dbStation, err := s.stationRepository.FindByID(ctx, station.ID, nil)
	if err != nil {
		return nil, err
	}

	var oldClusterId *uint
	if dbStation.ClusterID != nil {
		oldClusterId = dbStation.ClusterID
	}

	updatedStation, err := s.stationRepository.Update(ctx, station)
	if err != nil {
		return nil, err
	}

	if updatedStation.ClusterID != nil &&
		(oldClusterId == nil || *oldClusterId != *updatedStation.ClusterID) {

		event := events.StationEvent{
			Type:      events.StationAddedToCluster,
			StationId: updatedStation.ID,
			ClusterId: *updatedStation.ClusterID,
			Timestamp: time.Now(),
		}

		s.stationEventBus.Publish(&event)
	}

	// Station aus Cluster entfernt?
	if oldClusterId != nil &&
		(updatedStation.ClusterID == nil || *oldClusterId != *updatedStation.ClusterID) {
		event := events.StationEvent{
			Type:      events.StationRemovedFromCluster,
			StationId: updatedStation.ID,
			ClusterId: *oldClusterId,
			Timestamp: time.Now(),
		}

		s.stationEventBus.Publish(&event)
	}

	return updatedStation, nil
}

func (s *StationService) Create(ctx context.Context, station *models.Station, includeParam *string) (*models.Station, error) {
	result, err := s.stationRepository.Create(ctx, station)
	if err != nil {
		return nil, err
	}

	return result, nil
}
