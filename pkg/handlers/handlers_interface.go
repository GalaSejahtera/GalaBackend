package handlers

import (
	"context"
	pb "galasejahtera/pkg/api"

	"github.com/golang/protobuf/ptypes/empty"
)

// IHandlers ...
type IHandlers interface {
	// -------------- User ----------------
	CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error)
	GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error)
	GetPasswordReset(ctx context.Context, req *pb.GetPasswordResetRequest) (*pb.GetPasswordResetResponse, error)
	UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*empty.Empty, error)
	DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)
	DeleteUsers(ctx context.Context, req *pb.DeleteUsersRequest) (*pb.DeleteUsersResponse, error)
	UpdateUsers(ctx context.Context, req *pb.UpdateUsersRequest) (*pb.UpdateUsersResponse, error)
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	Logout(ctx context.Context, req *empty.Empty) (*empty.Empty, error)
	Refresh(ctx context.Context, req *empty.Empty) (*pb.RefreshResponse, error)
	GetUsersByZone(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error)
	GetRecentUsersByZone(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error)
	GetRecentUsersByUser(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error)
	// -------------- User ----------------

	// -------------- Zone ----------------
	GetRecentZonesByUser(ctx context.Context, req *pb.GetZonesRequest) (*pb.GetZonesResponse, error)
	CreateZone(ctx context.Context, req *pb.CreateZoneRequest) (*pb.CreateZoneResponse, error)
	GetZones(ctx context.Context, req *pb.GetZonesRequest) (*pb.GetZonesResponse, error)
	GetZone(ctx context.Context, req *pb.GetZoneRequest) (*pb.GetZoneResponse, error)
	DeleteZone(ctx context.Context, req *pb.DeleteZoneRequest) (*pb.DeleteZoneResponse, error)
	UpdateZone(ctx context.Context, req *pb.UpdateZoneRequest) (*pb.UpdateZoneResponse, error)
	DeleteZones(ctx context.Context, req *pb.DeleteZonesRequest) (*pb.DeleteZonesResponse, error)
	UpdateZones(ctx context.Context, req *pb.UpdateZonesRequest) (*pb.UpdateZonesResponse, error)
	// -------------- Zone ----------------

	// -------------- Faq ----------------
	CreateFaq(ctx context.Context, req *pb.CreateFaqRequest) (*pb.CreateFaqResponse, error)
	GetFaqs(ctx context.Context, req *pb.GetFaqsRequest) (*pb.GetFaqsResponse, error)
	GetFaq(ctx context.Context, req *pb.GetFaqRequest) (*pb.GetFaqResponse, error)
	DeleteFaq(ctx context.Context, req *pb.DeleteFaqRequest) (*pb.DeleteFaqResponse, error)
	UpdateFaq(ctx context.Context, req *pb.UpdateFaqRequest) (*pb.UpdateFaqResponse, error)
	DeleteFaqs(ctx context.Context, req *pb.DeleteFaqsRequest) (*pb.DeleteFaqsResponse, error)
	UpdateFaqs(ctx context.Context, req *pb.UpdateFaqsRequest) (*pb.UpdateFaqsResponse, error)
	// -------------- Zone ----------------

	// -------------- Client ----------------
	ClientGetNearbyUsers(ctx context.Context, req *pb.ClientGetNearbyUsersRequest) (*pb.ClientGetNearbyUsersResponse, error)
	ClientGetCurrentZones(ctx context.Context, req *pb.ClientGetCurrentZonesRequest) (*pb.ClientGetCurrentZonesResponse, error)
	ClientGetRecentZones(ctx context.Context, req *pb.ClientGetRecentZonesRequest) (*pb.GetZonesResponse, error)
	ClientCreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	ClientGetZones(ctx context.Context, req *pb.GetZonesRequest) (*pb.GetZonesResponse, error)
	ClientGetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error)
	ClientUpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error)
	// -------------- Client ----------------

	// -------------- Activity ----------------
	GetActivities(ctx context.Context, req *pb.GetActivitiesRequest) (*pb.GetActivitiesResponse, error)
	// -------------- Activity ----------------
}
