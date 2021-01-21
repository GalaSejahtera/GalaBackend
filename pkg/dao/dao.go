package dao

import (
	"context"

	"galasejahtera/pkg/dto"
)

// IUserDAO ...
type IUserDAO interface {
	// Create creates new user
	Create(ctx context.Context, user *dto.User) (*dto.User, error)
	// Update updates user
	Update(ctx context.Context, user *dto.User) (*dto.User, error)
	// Get gets user by ID
	Get(ctx context.Context, id string) (*dto.User, error)
	// BatchGet gets users by slice of IDs
	BatchGet(ctx context.Context, ids []string) ([]*dto.User, error)
	// Query queries users by sort, range, filter
	Query(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.User, error)
	// Delete deletes user by ID
	Delete(ctx context.Context, id string) error
	// BatchDelete deletes users by IDs
	BatchDelete(ctx context.Context, ids []string) ([]string, error)
	// GetNearbyUsers gets users within 50 meter
	GetNearbyUsers(ctx context.Context, user *dto.User) (int64, []*dto.User, error)
}

// IAuthDAO ...
type IAuthDAO interface {
	InitIndex(ctx context.Context) error
	// Create creates new auth token
	Create(ctx context.Context, auth *dto.AuthObject) (*dto.AuthObject, error)
	// Get gets auth token
	Get(ctx context.Context, token string) (*dto.AuthObject, error)
	// Delete deletes user by token
	Delete(ctx context.Context, token string) error
	// DeleteByID deletes user by ID
	DeleteByID(ctx context.Context, id string) error
}

// IReportDAO ...
type IReportDAO interface {
	// Create creates new report
	Create(ctx context.Context, report *dto.Report) (*dto.Report, error)
	// Get gets report
	Get(ctx context.Context, id string) (*dto.Report, error)
	// BatchGet gets reports by slice of IDs
	BatchGet(ctx context.Context, ids []string) ([]*dto.Report, error)
	// Query queries reports by sort, range, filter
	Query(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.Report, error)
	// Delete deletes report by ID
	Delete(ctx context.Context, id string) error
	// BatchDelete deletes reports by IDs
	BatchDelete(ctx context.Context, ids []string) ([]string, error)
	// Update updates report
	Update(ctx context.Context, report *dto.Report) (*dto.Report, error)
}
