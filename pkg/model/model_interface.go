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

	///////////// Faq models
	// CreateFaq creates new faq
	CreateFaq(ctx context.Context, Faq *dto.Faq) (*dto.Faq, error)
	// GetFaq gets activity by ID
	GetFaq(ctx context.Context, id string) (*dto.Faq, error)
	// BatchGetFaqs get faqs by slice of IDs
	BatchGetFaqs(ctx context.Context, ids []string) ([]*dto.Faq, error)
	// QueryFaqs queries faqs by sort, range, filter
	QueryFaqs(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.Faq, error)
	// DeleteFaq deletes faq by ID
	DeleteFaq(ctx context.Context, id string) (*dto.Faq, error)
	// DeleteFaqs delete faqs by IDs
	DeleteFaqs(ctx context.Context, ids []string) ([]string, error)
	// UpdateZone updates zone
	UpdateFaq(ctx context.Context, faq *dto.Faq) (*dto.Faq, error)
	// UpdateZones update faqs
	UpdateFaqs(ctx context.Context, faq *dto.Faq, ids []string) ([]string, error)
	/////////////
}
