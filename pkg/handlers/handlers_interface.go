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
	GetNearbyUsers(ctx context.Context, req *pb.GetNearbyUsersRequest) (*pb.GetNearbyUsersResponse, error)
	// -------------- User ----------------

	// -------------- Report ----------------
	CreateReport(ctx context.Context, req *pb.CreateReportRequest) (*pb.CreateReportResponse, error)
	GetReports(ctx context.Context, req *pb.GetReportsRequest) (*pb.GetReportsResponse, error)
	GetReport(ctx context.Context, req *pb.GetReportRequest) (*pb.GetReportResponse, error)
	DeleteReport(ctx context.Context, req *pb.DeleteReportRequest) (*pb.DeleteReportResponse, error)
	UpdateReport(ctx context.Context, req *pb.UpdateReportRequest) (*pb.UpdateReportResponse, error)
	DeleteReports(ctx context.Context, req *pb.DeleteReportsRequest) (*pb.DeleteReportsResponse, error)
	UpdateReports(ctx context.Context, req *pb.UpdateReportsRequest) (*pb.UpdateReportsResponse, error)
	// -------------- Report ----------------

	// -------------- Covid ----------------
	GetCovids(ctx context.Context, req *pb.GetCovidsRequest) (*pb.GetCovidsResponse, error)
	GetCovid(ctx context.Context, req *pb.GetCovidRequest) (*pb.GetCovidResponse, error)
	// -------------- Covid ----------------

	// -------------- Daily ----------------
	GetDistrict(ctx context.Context, req *pb.GetDistrictRequest) (*pb.GetDistrictResponse, error)
	// -------------- Daily ----------------

	GetKases(ctx context.Context, req *empty.Empty) (*pb.GetKasesResponse, error)
	GetRecentKases(ctx context.Context, req *empty.Empty) (*pb.GetRecentKasesResponse, error)
}
