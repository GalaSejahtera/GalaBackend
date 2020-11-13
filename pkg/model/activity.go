package model

import (
	"context"
	"safeworkout/pkg/dto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateActivity creates new activity
func (m *Model) CreateActivity(ctx context.Context, activity *dto.Activity) (*dto.Activity, error) {

	// check if Activity exist
	_, err := m.activityDAO.Get(ctx, activity.ID)

	// only can create activity if not found
	if err != nil && status.Code(err) == codes.Unknown {
		return m.activityDAO.Create(ctx, activity)
	}

	if err != nil {
		return nil, err
	}

	return nil, status.Error(codes.AlreadyExists, "Activity already exist!")
}

// GetActivity gets activity by ID
func (m *Model) GetActivity(ctx context.Context, id string) (*dto.Activity, error) {
	return m.activityDAO.Get(ctx, id)
}

// QueryActivities queries activities by sort, range, filter
func (m *Model) QueryActivities(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.Activity, error) {
	return m.activityDAO.Query(ctx, sort, itemsRange, filter)
}

// DeleteActivity deletes activity by ID
func (m *Model) DeleteActivity(ctx context.Context, id string) (*dto.Activity, error) {
	// check if activity exist
	u, err := m.activityDAO.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = m.activityDAO.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// DeleteActivities delete activities by IDs
func (m *Model) DeleteActivities(ctx context.Context, ids []string) ([]string, error) {
	var deletedIDs []string
	for _, id := range ids {
		u, err := m.DeleteActivity(ctx, id)
		if err != nil {
			return nil, err
		}
		deletedIDs = append(deletedIDs, u.ID)
	}

	return deletedIDs, nil
}
