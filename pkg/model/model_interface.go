package model

import (
	"context"
	"galasejahtera/pkg/dto"
)

// IModel ...
type IModel interface {
	///////////// User models
	// CreateUser creates new user
	CreateUser(ctx context.Context, user *dto.User) (*dto.User, error)
	// UpdateUser updates user
	UpdateUser(ctx context.Context, user *dto.User) (*dto.User, error)
	// UpdateUserPassword updates user password only
	UpdateUserPassword(ctx context.Context, user *dto.User) (*dto.User, error)
	// CreateToken creates token with custom ttl
	CreateToken(ctx context.Context, auth *dto.AuthObject) (*dto.AuthObject, error)
	// RevokeTokensByUserID revokes all tokens by UserID
	RevokeTokensByUserID(ctx context.Context, id string) error
	// GetUserIDByToken gets userID by token
	GetUserIDByToken(ctx context.Context, token string) (string, error)
	// UpdateUsers update users
	UpdateUsers(ctx context.Context, user *dto.User, ids []string) ([]string, error)
	// GetUser gets user by ID
	GetUser(ctx context.Context, id string) (*dto.User, error)
	// BatchGetUsers get users by slice of IDs
	BatchGetUsers(ctx context.Context, ids []string) ([]*dto.User, error)
	// QueryUsers queries users by sort, range, filter
	QueryUsers(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.User, error)
	// DeleteUser deletes user by ID
	DeleteUser(ctx context.Context, id string) (*dto.User, error)
	// RevokeUserTokens revoke all user tokens
	RevokeUserTokens(ctx context.Context) error
	// DeleteUsers delete users by IDs
	DeleteUsers(ctx context.Context, ids []string) ([]string, error)
	// Login verifies user by email and password and return tokens
	Login(ctx context.Context, email string, password string) (*dto.User, error)
	// VerifyUser verifies user by header
	VerifyUser(ctx context.Context, header string) (*dto.User, error)
	// Logout logs user out from the system by header
	Logout(ctx context.Context, header string) error
	// Refresh returns new token to authorized user by header
	Refresh(ctx context.Context, header string) (*dto.User, error)
	// GetNearbyUsers get nearby users count given user
	GetNearbyUsers(ctx context.Context, user *dto.User) (int64, []*dto.User, error)
	// DisableInactiveUsers disable inactive users
	DisableInactiveUsers(ctx context.Context) error
	/////////////

	///////////// Report models
	// CreateReport creates new report
	CreateReport(ctx context.Context, Report *dto.Report) (*dto.Report, error)
	// GetReport gets activity by ID
	GetReport(ctx context.Context, id string) (*dto.Report, error)
	// BatchGetReports get reports by slice of IDs
	BatchGetReports(ctx context.Context, ids []string) ([]*dto.Report, error)
	// QueryReports queries reports by sort, range, filter
	QueryReports(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.Report, error)
	// DeleteReport deletes report by ID
	DeleteReport(ctx context.Context, id string) (*dto.Report, error)
	// DeleteReports delete reports by IDs
	DeleteReports(ctx context.Context, ids []string) ([]string, error)
	// UpdateZone updates zone
	UpdateReport(ctx context.Context, report *dto.Report) (*dto.Report, error)
	// UpdateZones update reports
	UpdateReports(ctx context.Context, report *dto.Report, ids []string) ([]string, error)
	/////////////

	///////////// Covid models
	// GetCovid gets activity by ID
	GetCovid(ctx context.Context, id string) (*dto.Covid, error)
	// QueryCovids queries covids by sort, range, filter
	QueryCovids(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.Covid, error)
	/////////////
}
