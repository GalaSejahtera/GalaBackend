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
	// QueryRecentZonesByUserID queries past 14 days zones by userID
	QueryRecentZonesByUserID(ctx context.Context, userID string) ([]*dto.Zone, error)
	// QueryRecentUsersByUserID queries past 14 days users by userID
	QueryRecentUsersByUserID(ctx context.Context, userID string) ([]*dto.User, error)
}

// IZoneDAO ...
type IZoneDAO interface {
	// Create creates new zone
	Create(ctx context.Context, zone *dto.Zone) (*dto.Zone, error)
	// Update updates zone
	Update(ctx context.Context, zone *dto.Zone) (*dto.Zone, error)
	// Get gets zone by ID
	Get(ctx context.Context, id string) (*dto.Zone, error)
	// BatchGet gets zones by slice of IDs
	BatchGet(ctx context.Context, ids []string) ([]*dto.Zone, error)
	// Query queries zones by sort, range, filter
	Query(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.Zone, error)
	// Delete deletes zone by ID
	Delete(ctx context.Context, id string) error
	// BatchDelete deletes zones by IDs
	BatchDelete(ctx context.Context, ids []string) ([]string, error)
	// GetByUser gets zone and sub zones given user
	GetByUser(ctx context.Context, user *dto.User) (*dto.Zone, []*dto.Zone, error)
	// GetUsersByZone gets users by zone
	GetUsersByZone(ctx context.Context, zone *dto.Zone) ([]*dto.User, error)
	// GetSubZones gets sub zones by zone
	GetSubZones(ctx context.Context, zone *dto.Zone, user *dto.User) ([]*dto.Zone, error)
	// QueryRecentUsersByZoneID queries past 14 days users by zoneID
	QueryRecentUsersByZoneID(ctx context.Context, zoneID string) ([]*dto.User, error)
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

// IActivityDAO ...
type IActivityDAO interface {
	// Create creates new activity
	Create(ctx context.Context, activity *dto.Activity) (*dto.Activity, error)
	// Get gets activity
	Get(ctx context.Context, id string) (*dto.Activity, error)
	// Query queries activities by sort, range, filter
	Query(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.Activity, error)
	// Delete deletes activity by ID
	Delete(ctx context.Context, id string) error
	// BatchDelete deletes activities by IDs
	BatchDelete(ctx context.Context, ids []string) ([]string, error)
}

// IFaqDAO ...
type IFaqDAO interface {
	// Create creates new faq
	Create(ctx context.Context, faq *dto.Faq) (*dto.Faq, error)
	// Get gets faq
	Get(ctx context.Context, id string) (*dto.Faq, error)
	// BatchGet gets faqs by slice of IDs
	BatchGet(ctx context.Context, ids []string) ([]*dto.Faq, error)
	// Query queries faqs by sort, range, filter
	Query(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.Faq, error)
	// Delete deletes faq by ID
	Delete(ctx context.Context, id string) error
	// BatchDelete deletes faqs by IDs
	BatchDelete(ctx context.Context, ids []string) ([]string, error)
	// Update updates faq
	Update(ctx context.Context, faq *dto.Faq) (*dto.Faq, error)
}
