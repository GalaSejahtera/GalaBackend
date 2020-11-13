package model

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/dto"
)

// CreateZone creates new zone
func (m *Model) CreateZone(ctx context.Context, zone *dto.Zone) (*dto.Zone, error) {

	// check if Zone exist
	_, err := m.zoneDAO.Get(ctx, zone.ID)

	// only can create zone if not found
	if err != nil && status.Code(err) == codes.Unknown {
		return m.zoneDAO.Create(ctx, zone)
	}

	if err != nil {
		return nil, err
	}

	return nil, status.Error(codes.AlreadyExists, "Zone already exist!")
}

// UpdateZone updates zone
func (m *Model) UpdateZone(ctx context.Context, zone *dto.Zone) (*dto.Zone, error) {

	// check if zone exists
	u, err := m.zoneDAO.Get(ctx, zone.ID)
	if err != nil {
		return nil, err
	}

	// patch zone
	u.Name = zone.Name
	u.Lat = zone.Lat
	u.Long = zone.Long
	u.Type = zone.Type
	u.Radius = zone.Radius
	u.Capacity = zone.Capacity

	_, err = m.zoneDAO.Update(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// UpdateZones update zones
func (m *Model) UpdateZones(ctx context.Context, zone *dto.Zone, ids []string) ([]string, error) {
	if len(ids) > 1 {
		return nil, constants.OperationUnsupportedError
	}
	zone.ID = ids[0]
	u, err := m.UpdateZone(ctx, zone)
	if err != nil {
		return nil, err
	}
	return []string{u.ID}, err
}

// GetZone gets zone by ID
func (m *Model) GetZone(ctx context.Context, id string) (*dto.Zone, error) {
	return m.zoneDAO.Get(ctx, id)
}

// BatchGetZones get zones by slice of IDs
func (m *Model) BatchGetZones(ctx context.Context, ids []string) ([]*dto.Zone, error) {
	return m.zoneDAO.BatchGet(ctx, ids)
}

// QueryZones queries zones by sort, range, filter
func (m *Model) QueryZones(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.Zone, error) {
	return m.zoneDAO.Query(ctx, sort, itemsRange, filter)
}

// DeleteZone deletes zone by ID
func (m *Model) DeleteZone(ctx context.Context, id string) (*dto.Zone, error) {
	// check if zone exist
	u, err := m.zoneDAO.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = m.zoneDAO.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// DeleteZones delete zones by IDs
func (m *Model) DeleteZones(ctx context.Context, ids []string) ([]string, error) {
	var deletedIDs []string
	for _, id := range ids {
		u, err := m.DeleteZone(ctx, id)
		if err != nil {
			return nil, err
		}
		deletedIDs = append(deletedIDs, u.ID)
	}

	return deletedIDs, nil
}

// GetZonesByUser gets zones and sub zones given userID
func (m *Model) GetZonesByUser(ctx context.Context, userID string) (*dto.Zone, []*dto.Zone, error) {
	user, err := m.userDAO.Get(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	zone, subZones, err := m.zoneDAO.GetByUser(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	return zone, subZones, err
}

// GetRecentZonesByUser gets recent zones given userID
func (m *Model) GetRecentZonesByUser(ctx context.Context, userID string) ([]*dto.Zone, error) {
	return m.userDAO.QueryRecentZonesByUserID(ctx, userID)
}
