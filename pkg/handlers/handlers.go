package handlers

import (
	"context"
	"errors"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/handlers/covid"
	"galasejahtera/pkg/handlers/daily"
	"galasejahtera/pkg/handlers/kase"
	"galasejahtera/pkg/handlers/report"
	"galasejahtera/pkg/handlers/user"
	"galasejahtera/pkg/logger"
	"galasejahtera/pkg/model"
	"os"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

// Handlers ...
type Handlers struct {
	Model model.IModel
}

// NewHandlers ...
func NewHandlers(model model.IModel) IHandlers {
	return &Handlers{Model: model}
}

func (s *Handlers) GetKases(ctx context.Context, req *empty.Empty) (*pb.GetKasesResponse, error) {
	handler := &kase.GetKasesHandler{Model: s.Model}
	resp, err := handler.GetKases(ctx, req)
	if err != nil {
		logger.Log.Error("GetKasesHandler: " + err.Error())
		return nil, err
	}
	return resp, nil
}

func (s *Handlers) GetRecentKases(ctx context.Context, req *empty.Empty) (*pb.GetRecentKasesResponse, error) {
	handler := &kase.GetRecentKasesHandler{Model: s.Model}
	resp, err := handler.GetRecentKases(ctx, req)
	if err != nil {
		logger.Log.Error("GetRecentKasesHandler: " + err.Error())
		return nil, err
	}
	return resp, nil
}

func (s *Handlers) GetNearbyUsers(ctx context.Context, req *pb.GetNearbyUsersRequest) (*pb.GetNearbyUsersResponse, error) {
	handler := &user.GetNearbyUsersHandler{Model: s.Model}
	resp, err := handler.GetNearbyUsers(ctx, req)
	if err != nil {
		logger.Log.Error("ClientGetNearbyUsersHandler: " + err.Error())
		return nil, err
	}
	return resp, nil
}

func (s *Handlers) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	handler := &user.CreateUserHandler{Model: s.Model}
	resp, err := handler.CreateUser(ctx, req)
	if err != nil {
		logger.Log.Error("CreateUserHandler: " + err.Error())
		return nil, err
	}
	return resp, nil
}

func (s *Handlers) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	handler := &user.GetUsersHandler{Model: s.Model}
	resp, err := handler.GetUsers(ctx, req)
	if err != nil {
		logger.Log.Error("GetUsersHandler: " + err.Error())
		return nil, err
	}
	return resp, nil
}

func (s *Handlers) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	handler := &user.GetUserHandler{Model: s.Model}
	resp, err := handler.GetUser(ctx, req)
	if err != nil {
		logger.Log.Error("GetUserHandler: "+err.Error(), zap.String("UserID", req.Id))
		return nil, err
	}
	logger.Log.Info("GetUserHandler", zap.String("UserID", req.Id))
	return resp, nil
}

func (s *Handlers) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	handler := &user.DeleteUserHandler{Model: s.Model}
	resp, err := handler.DeleteUser(ctx, req)
	if err != nil {
		logger.Log.Error("DeleteUserHandler: "+err.Error(), zap.String("TargetUserID", req.Id))
		return nil, err
	}
	logger.Log.Info("DeleteUserHandler", zap.String("TargetUserID", req.Id))
	return resp, nil
}

func (s *Handlers) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	handler := &user.UpdateUserHandler{Model: s.Model}
	resp, err := handler.UpdateUser(ctx, req)
	if err != nil {
		logger.Log.Error("UpdateUserHandler: "+err.Error(), zap.String("TargetUserID", req.Id))
		return nil, err
	}
	logger.Log.Info("UpdateUserHandler", zap.String("TargetUserID", req.Id))
	return resp, nil
}

func (s *Handlers) DeleteUsers(ctx context.Context, req *pb.DeleteUsersRequest) (*pb.DeleteUsersResponse, error) {
	handler := &user.DeleteUsersHandler{Model: s.Model}
	resp, err := handler.DeleteUsers(ctx, req)
	if err != nil {
		logger.Log.Error("DeleteUsersHandler: "+err.Error(), zap.Strings("TargetUserIDs", req.Ids))
		return nil, err
	}
	logger.Log.Info("DeleteUsersHandler", zap.Strings("TargetUserIDs", req.Ids))
	return resp, nil
}

func (s *Handlers) UpdateUsers(ctx context.Context, req *pb.UpdateUsersRequest) (*pb.UpdateUsersResponse, error) {
	handler := &user.UpdateUsersHandler{Model: s.Model}
	resp, err := handler.UpdateUsers(ctx, req)
	if err != nil {
		logger.Log.Error("UpdateUsersHandler: "+err.Error(), zap.Strings("TargetUserIDs", req.Ids))
		return nil, err
	}
	logger.Log.Info("UpdateUsersHandler", zap.Strings("TargetUserIDs", req.Ids))
	return resp, nil
}

func (s *Handlers) GetPasswordReset(ctx context.Context, req *pb.GetPasswordResetRequest) (*pb.GetPasswordResetResponse, error) {
	handler := &user.GetPasswordResetHandler{Model: s.Model}
	resp, err := handler.GetPasswordReset(ctx, req)
	if err != nil {
		logger.Log.Error("GetPasswordResetHandler: "+err.Error(), zap.String("TargetUserID", req.Id))
		return nil, err
	}
	logger.Log.Info("GetPasswordResetHandler", zap.String("TargetUserID", req.Id))
	return resp, nil
}

func (s *Handlers) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*empty.Empty, error) {
	handler := &user.UpdatePasswordHandler{Model: s.Model}
	resp, err := handler.UpdatePassword(ctx, req)
	if err != nil {
		logger.Log.Error("UpdatePasswordHandler: " + err.Error())
		return nil, err
	}
	logger.Log.Info("UpdatePasswordHandler")
	return resp, nil
}

func (s *Handlers) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	handler := &user.LoginHandler{Model: s.Model}
	resp, err := handler.Login(ctx, req)
	if err != nil {
		logger.Log.Error("LoginHandler: "+err.Error(), zap.String("email", req.Email))
		return nil, err
	}
	logger.Log.Info("LoginHandler", zap.String("email", req.Email))
	return resp, nil
}

func (s *Handlers) Logout(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	handler := &user.LogoutHandler{Model: s.Model}
	resp, err := handler.Logout(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *Handlers) Refresh(ctx context.Context, _ *empty.Empty) (*pb.RefreshResponse, error) {
	handler := &user.RefreshHandler{Model: s.Model}
	resp, err := handler.Refresh(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// -------------------- Covids ------------------------

func (s *Handlers) GetCovids(ctx context.Context, req *pb.GetCovidsRequest) (*pb.GetCovidsResponse, error) {
	handler := &covid.GetCovidsHandler{Model: s.Model}
	resp, err := handler.GetCovids(ctx, req)
	if err != nil {
		logger.Log.Error("GetCovidsHandler: " + err.Error())
		return nil, err
	}
	return resp, nil
}

func (s *Handlers) GetCovid(ctx context.Context, req *pb.GetCovidRequest) (*pb.GetCovidResponse, error) {
	handler := &covid.GetCovidHandler{Model: s.Model}
	resp, err := handler.GetCovid(ctx, req)
	if err != nil {
		logger.Log.Error("GetCovidHandler: "+err.Error(), zap.String("CovidID", req.Id))
		return nil, err
	}
	logger.Log.Info("GetCovidHandler", zap.String("CovidID", req.Id))
	return resp, nil
}

// -------------------- Reports ------------------------

func (s *Handlers) CreateReport(ctx context.Context, req *pb.CreateReportRequest) (*pb.CreateReportResponse, error) {
	handler := &report.CreateReportHandler{Model: s.Model}
	resp, err := handler.CreateReport(ctx, req)
	if err != nil {
		logger.Log.Error("CreateReportHandler: " + err.Error())
		return nil, err
	}
	logger.Log.Info("CreateReportHandler")
	return resp, nil
}

func (s *Handlers) GetReports(ctx context.Context, req *pb.GetReportsRequest) (*pb.GetReportsResponse, error) {
	handler := &report.GetReportsHandler{Model: s.Model}
	resp, err := handler.GetReports(ctx, req)
	if err != nil {
		logger.Log.Error("GetReportsHandler: " + err.Error())
		return nil, err
	}
	return resp, nil
}

func (s *Handlers) GetReport(ctx context.Context, req *pb.GetReportRequest) (*pb.GetReportResponse, error) {
	u, err := s.validateUser(ctx, constants.AllCanAccess)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &report.GetReportHandler{Model: s.Model}
	resp, err := handler.GetReport(ctx, req)
	if err != nil {
		logger.Log.Error("GetReportHandler: "+err.Error(), zap.String("UserID", u.ID), zap.String("ReportID", req.Id))
		return nil, err
	}
	logger.Log.Info("GetReportHandler", zap.String("UserID", u.ID), zap.String("ReportID", req.Id))
	return resp, nil
}

func (s *Handlers) DeleteReport(ctx context.Context, req *pb.DeleteReportRequest) (*pb.DeleteReportResponse, error) {
	u, err := s.validateUser(ctx, constants.AllCanAccess)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &report.DeleteReportHandler{Model: s.Model}
	resp, err := handler.DeleteReport(ctx, req)
	if err != nil {
		logger.Log.Error("DeleteReportHandler: "+err.Error(), zap.String("UserID", u.ID), zap.String("TargetReportID", req.Id))
		return nil, err
	}
	logger.Log.Info("DeleteReportHandler", zap.String("UserID", u.ID), zap.String("TargetReportID", req.Id))
	return resp, nil
}

func (s *Handlers) UpdateReport(ctx context.Context, req *pb.UpdateReportRequest) (*pb.UpdateReportResponse, error) {
	u, err := s.validateUser(ctx, constants.AllCanAccess)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &report.UpdateReportHandler{Model: s.Model}
	resp, err := handler.UpdateReport(ctx, req)
	if err != nil {
		logger.Log.Error("UpdateReportHandler: "+err.Error(), zap.String("UserID", u.ID), zap.String("TargetReportID", req.Id))
		return nil, err
	}
	logger.Log.Info("UpdateReportHandler", zap.String("UserID", u.ID), zap.String("TargetReportID", req.Id))
	return resp, nil
}

func (s *Handlers) DeleteReports(ctx context.Context, req *pb.DeleteReportsRequest) (*pb.DeleteReportsResponse, error) {
	u, err := s.validateUser(ctx, constants.AllCanAccess)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &report.DeleteReportsHandler{Model: s.Model}
	resp, err := handler.DeleteReports(ctx, req)
	if err != nil {
		logger.Log.Error("DeleteReportsHandler: "+err.Error(), zap.String("UserID", u.ID), zap.Strings("TargetReportIDs", req.Ids))
		return nil, err
	}
	logger.Log.Info("DeleteReportsHandler", zap.String("UserID", u.ID), zap.Strings("TargetReportIDs", req.Ids))
	return resp, nil
}

func (s *Handlers) UpdateReports(ctx context.Context, req *pb.UpdateReportsRequest) (*pb.UpdateReportsResponse, error) {
	u, err := s.validateUser(ctx, constants.AllCanAccess)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &report.UpdateReportsHandler{Model: s.Model}
	resp, err := handler.UpdateReports(ctx, req)
	if err != nil {
		logger.Log.Error("UpdateReportsHandler: "+err.Error(), zap.String("UserID", u.ID), zap.Strings("TargetReportIDs", req.Ids))
		return nil, err
	}
	logger.Log.Info("UpdateReportsHandler", zap.String("UserID", u.ID), zap.Strings("TargetReportIDs", req.Ids))
	return resp, nil
}

// -------------------- Reports ------------------------

// -------------------- Daily ------------------------

func (s *Handlers) GetDistrict(ctx context.Context, req *pb.GetDistrictRequest) (*pb.GetDistrictResponse, error) {
	handler := &daily.GetDailyHandler{Model: s.Model}
	resp, err := handler.GetDistrict(ctx, req)
	if err != nil {
		logger.Log.Error("GetDistrictHandler: "+err.Error(), zap.String("District", req.Id))
		return nil, err
	}
	logger.Log.Info("GetDistrictHandler", zap.String("District", req.Id))
	return resp, nil
}

// -------------------- Daily ------------------------

func (s *Handlers) validateUser(ctx context.Context, roles []string) (*dto.User, error) {
	if os.Getenv("AUTH_ENABLED") != "true" {
		return &dto.User{}, nil
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("ValidateUser: metadata not found")
	}
	tokenSlice := md.Get("authorization")
	if len(tokenSlice) < 1 {
		return nil, errors.New("ValidateUser: token not found")
	}
	token := tokenSlice[0]

	// exemption: backend user
	if token == os.Getenv("BACKEND_USER") {
		// get user from collection
		u, err := s.Model.GetUser(ctx, token)
		if err != nil {
			return &dto.User{}, nil
		}
		return u, nil
	}

	// new verify user mechanism
	u, err := s.Model.VerifyUser(ctx, strings.Join(tokenSlice, " "))
	if err != nil {
		return nil, err
	}

	// check if user is allowed to access the API
	authorized := false
	for _, role := range roles {
		if u.Role == role {
			authorized = true
		}
	}
	if !authorized {
		return nil, errors.New("unauthorized access")
	}
	return u, nil
}

// -------------------- Users ------------------------
