package handlers

import (
	"context"
	"errors"
	"os"
	pb "safeworkout/pkg/api"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/dto"
	"safeworkout/pkg/handlers/activity"
	"safeworkout/pkg/handlers/client"
	"safeworkout/pkg/handlers/faq"
	"safeworkout/pkg/handlers/user"
	"safeworkout/pkg/handlers/zone"
	"safeworkout/pkg/logger"
	"safeworkout/pkg/model"
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

// -------------------- Client ------------------------

func (s *Handlers) ClientGetRecentZones(ctx context.Context, req *pb.ClientGetRecentZonesRequest) (*pb.GetZonesResponse, error) {
	handler := &client.GetRecentZonesHandler{Model: s.Model}
	resp, err := handler.GetRecentZones(ctx, req)
	if err != nil {
		logger.Log.Error("GetRecentZonesClientHandler: " + err.Error())
		return nil, err
	}
	logger.Log.Info("GetRecentZonesClientHandler")
	return resp, nil
}

func (s *Handlers) ClientCreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	handler := &client.CreateUserClientHandler{Model: s.Model}
	resp, err := handler.CreateUser(ctx, req)
	if err != nil {
		logger.Log.Error("CreateUserClientHandler: " + err.Error())
		return nil, err
	}
	logger.Log.Info("CreateUserClientHandler")
	return resp, nil
}

func (s *Handlers) ClientUpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	u, err := s.validateUser(ctx, constants.AllCanAccess)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &client.UpdateUserClientHandler{Model: s.Model}
	resp, err := handler.UpdateUser(ctx, req)
	if err != nil {
		logger.Log.Error("ClientUpdateUserHandler: "+err.Error(), zap.String("UserID", u.ID), zap.String("TargetUserID", req.Id))
		return nil, err
	}
	logger.Log.Info("ClientUpdateUserHandler", zap.String("UserID", u.ID), zap.String("TargetUserID", req.Id))
	return resp, nil
}

func (s *Handlers) ClientGetNearbyUsers(ctx context.Context, req *pb.ClientGetNearbyUsersRequest) (*pb.ClientGetNearbyUsersResponse, error) {
	u, err := s.validateUser(ctx, constants.AllCanAccess)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &client.GetNearbyUsersHandler{Model: s.Model}
	resp, err := handler.GetNearbyUsers(ctx, req)
	if err != nil {
		logger.Log.Error("ClientGetNearbyUsersHandler: "+err.Error(), zap.String("UserID", u.ID))
		return nil, err
	}
	logger.Log.Info("ClientGetNearbyUsersHandler", zap.String("UserID", u.ID))
	return resp, nil
}

func (s *Handlers) ClientGetCurrentZones(ctx context.Context, req *pb.ClientGetCurrentZonesRequest) (*pb.ClientGetCurrentZonesResponse, error) {
	u, err := s.validateUser(ctx, constants.AllCanAccess)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &client.GetCurrentZonesHandler{Model: s.Model}
	resp, err := handler.GetCurrentZones(ctx, req)
	if err != nil {
		logger.Log.Error("ClientGetCurrentZonesHandler: "+err.Error(), zap.String("UserID", u.ID))
		return nil, err
	}
	logger.Log.Info("ClientGetCurrentZonesHandler", zap.String("UserID", u.ID))
	return resp, nil
}

func (s *Handlers) ClientGetZones(ctx context.Context, req *pb.GetZonesRequest) (*pb.GetZonesResponse, error) {
	u, err := s.validateUser(ctx, constants.AllCanAccess)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &zone.GetZonesHandler{Model: s.Model}
	resp, err := handler.GetZones(ctx, req)
	if err != nil {
		logger.Log.Error("GetZonesClientHandler: "+err.Error(), zap.String("UserID", u.ID))
		return nil, err
	}
	logger.Log.Info("GetZonesClientHandler", zap.String("UserID", u.ID))
	return resp, nil
}
func (s *Handlers) ClientGetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	u, err := s.validateUser(ctx, constants.AllCanAccess)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &client.GetUsersClientHandler{Model: s.Model}
	resp, err := handler.GetUsers(ctx, req)
	if err != nil {
		logger.Log.Error("GetUsersClientHandler: "+err.Error(), zap.String("UserID", u.ID))
		return nil, err
	}
	logger.Log.Info("GetUsersClientHandler", zap.String("UserID", u.ID))
	return resp, nil
}

// -------------------- Client ------------------------

// -------------------- Zones ------------------------

func (s *Handlers) CreateZone(ctx context.Context, req *pb.CreateZoneRequest) (*pb.CreateZoneResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserAndAdmin)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &zone.CreateZoneHandler{Model: s.Model}
	resp, err := handler.CreateZone(ctx, req)
	if err != nil {
		logger.Log.Error("CreateZoneHandler: "+err.Error(), zap.String("UserID", u.ID))
		return nil, err
	}
	logger.Log.Info("CreateZoneHandler", zap.String("UserID", u.ID))
	return resp, nil
}

func (s *Handlers) GetZones(ctx context.Context, req *pb.GetZonesRequest) (*pb.GetZonesResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserAndAdmin)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &zone.GetZonesHandler{Model: s.Model}
	resp, err := handler.GetZones(ctx, req)
	if err != nil {
		logger.Log.Error("GetZonesHandler: "+err.Error(), zap.String("UserID", u.ID))
		return nil, err
	}
	logger.Log.Info("GetZonesHandler", zap.String("UserID", u.ID))
	return resp, nil
}

func (s *Handlers) GetZone(ctx context.Context, req *pb.GetZoneRequest) (*pb.GetZoneResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserAndAdmin)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &zone.GetZoneHandler{Model: s.Model}
	resp, err := handler.GetZone(ctx, req)
	if err != nil {
		logger.Log.Error("GetZoneHandler: "+err.Error(), zap.String("UserID", u.ID), zap.String("ZoneID", req.Id))
		return nil, err
	}
	logger.Log.Info("GetZoneHandler", zap.String("UserID", u.ID), zap.String("ZoneID", req.Id))
	return resp, nil
}

func (s *Handlers) DeleteZone(ctx context.Context, req *pb.DeleteZoneRequest) (*pb.DeleteZoneResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserOnly)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &zone.DeleteZoneHandler{Model: s.Model}
	resp, err := handler.DeleteZone(ctx, req)
	if err != nil {
		logger.Log.Error("DeleteZoneHandler: "+err.Error(), zap.String("UserID", u.ID), zap.String("TargetZoneID", req.Id))
		return nil, err
	}
	logger.Log.Info("DeleteZoneHandler", zap.String("UserID", u.ID), zap.String("TargetZoneID", req.Id))
	return resp, nil
}

func (s *Handlers) UpdateZone(ctx context.Context, req *pb.UpdateZoneRequest) (*pb.UpdateZoneResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserAndAdmin)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &zone.UpdateZoneHandler{Model: s.Model}
	resp, err := handler.UpdateZone(ctx, req)
	if err != nil {
		logger.Log.Error("UpdateZoneHandler: "+err.Error(), zap.String("UserID", u.ID), zap.String("TargetZoneID", req.Id))
		return nil, err
	}
	logger.Log.Info("UpdateZoneHandler", zap.String("UserID", u.ID), zap.String("TargetZoneID", req.Id))
	return resp, nil
}

func (s *Handlers) DeleteZones(ctx context.Context, req *pb.DeleteZonesRequest) (*pb.DeleteZonesResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserOnly)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &zone.DeleteZonesHandler{Model: s.Model}
	resp, err := handler.DeleteZones(ctx, req)
	if err != nil {
		logger.Log.Error("DeleteZonesHandler: "+err.Error(), zap.String("UserID", u.ID), zap.Strings("TargetZoneIDs", req.Ids))
		return nil, err
	}
	logger.Log.Info("DeleteZonesHandler", zap.String("UserID", u.ID), zap.Strings("TargetZoneIDs", req.Ids))
	return resp, nil
}

func (s *Handlers) UpdateZones(ctx context.Context, req *pb.UpdateZonesRequest) (*pb.UpdateZonesResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserAndAdmin)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &zone.UpdateZonesHandler{Model: s.Model}
	resp, err := handler.UpdateZones(ctx, req)
	if err != nil {
		logger.Log.Error("UpdateZonesHandler: "+err.Error(), zap.String("UserID", u.ID), zap.Strings("TargetZoneIDs", req.Ids))
		return nil, err
	}
	logger.Log.Info("UpdateZonesHandler", zap.String("UserID", u.ID), zap.Strings("TargetZoneIDs", req.Ids))
	return resp, nil
}

// -------------------- Zones ------------------------

func (s *Handlers) GetRecentUsersByZone(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserAndAdmin)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &user.GetRecentUsersByZoneHandler{Model: s.Model}
	resp, err := handler.GetRecentUsersByZone(ctx, req)
	if err != nil {
		logger.Log.Error("GetRecentUsersByZoneHandler: "+err.Error(), zap.String("UserID", u.ID))
		return nil, err
	}
	logger.Log.Info("GetRecentUsersByZoneHandler", zap.String("UserID", u.ID))
	return resp, nil
}

func (s *Handlers) GetRecentUsersByUser(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserAndAdmin)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &user.GetRecentUsersByUserHandler{Model: s.Model}
	resp, err := handler.GetRecentUsersByUser(ctx, req)
	if err != nil {
		logger.Log.Error("GetRecentUsersByUserHandler: "+err.Error(), zap.String("UserID", u.ID))
		return nil, err
	}
	logger.Log.Info("GetRecentUsersByUserHandler", zap.String("UserID", u.ID))
	return resp, nil
}

func (s *Handlers) GetRecentZonesByUser(ctx context.Context, req *pb.GetZonesRequest) (*pb.GetZonesResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserAndAdmin)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &zone.GetRecentZonesByUserHandler{Model: s.Model}
	resp, err := handler.GetRecentZonesByUser(ctx, req)
	if err != nil {
		logger.Log.Error("GetRecentZonesByUserHandler: "+err.Error(), zap.String("UserID", u.ID))
		return nil, err
	}
	logger.Log.Info("GetRecentZonesByUserHandler", zap.String("UserID", u.ID))
	return resp, nil
}

func (s *Handlers) GetUsersByZone(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserAndAdmin)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &user.GetUsersByZoneHandler{Model: s.Model}
	resp, err := handler.GetUsersByZone(ctx, req)
	if err != nil {
		logger.Log.Error("GetUsersByZoneHandler: "+err.Error(), zap.String("UserID", u.ID))
		return nil, err
	}
	logger.Log.Info("GetUsersByZoneHandler", zap.String("UserID", u.ID))
	return resp, nil
}

func (s *Handlers) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserOnly)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &user.CreateUserHandler{Model: s.Model}
	resp, err := handler.CreateUser(ctx, req)
	if err != nil {
		logger.Log.Error("CreateUserHandler: "+err.Error(), zap.String("UserID", u.ID))
		return nil, err
	}
	logger.Log.Info("CreateUserHandler", zap.String("UserID", u.ID))
	return resp, nil
}

func (s *Handlers) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserAndAdmin)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &user.GetUsersHandler{Model: s.Model}
	resp, err := handler.GetUsers(ctx, req)
	if err != nil {
		logger.Log.Error("GetUsersHandler: "+err.Error(), zap.String("UserID", u.ID))
		return nil, err
	}
	logger.Log.Info("GetUsersHandler", zap.String("UserID", u.ID))
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
	u, err := s.validateUser(ctx, constants.SuperUserOnly)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &user.DeleteUserHandler{Model: s.Model}
	resp, err := handler.DeleteUser(ctx, req)
	if err != nil {
		logger.Log.Error("DeleteUserHandler: "+err.Error(), zap.String("UserID", u.ID), zap.String("TargetUserID", req.Id))
		return nil, err
	}
	logger.Log.Info("DeleteUserHandler", zap.String("UserID", u.ID), zap.String("TargetUserID", req.Id))
	return resp, nil
}

func (s *Handlers) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserOnly)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &user.UpdateUserHandler{Model: s.Model}
	resp, err := handler.UpdateUser(ctx, req)
	if err != nil {
		logger.Log.Error("UpdateUserHandler: "+err.Error(), zap.String("UserID", u.ID), zap.String("TargetUserID", req.Id))
		return nil, err
	}
	logger.Log.Info("UpdateUserHandler", zap.String("UserID", u.ID), zap.String("TargetUserID", req.Id))
	return resp, nil
}

func (s *Handlers) DeleteUsers(ctx context.Context, req *pb.DeleteUsersRequest) (*pb.DeleteUsersResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserOnly)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &user.DeleteUsersHandler{Model: s.Model}
	resp, err := handler.DeleteUsers(ctx, req)
	if err != nil {
		logger.Log.Error("DeleteUsersHandler: "+err.Error(), zap.String("UserID", u.ID), zap.Strings("TargetUserIDs", req.Ids))
		return nil, err
	}
	logger.Log.Info("DeleteUsersHandler", zap.String("UserID", u.ID), zap.Strings("TargetUserIDs", req.Ids))
	return resp, nil
}

func (s *Handlers) UpdateUsers(ctx context.Context, req *pb.UpdateUsersRequest) (*pb.UpdateUsersResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserOnly)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &user.UpdateUsersHandler{Model: s.Model}
	resp, err := handler.UpdateUsers(ctx, req)
	if err != nil {
		logger.Log.Error("UpdateUsersHandler: "+err.Error(), zap.String("UserID", u.ID), zap.Strings("TargetUserIDs", req.Ids))
		return nil, err
	}
	logger.Log.Info("UpdateUsersHandler", zap.String("UserID", u.ID), zap.Strings("TargetUserIDs", req.Ids))
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

// -------------------- Faqs ------------------------

func (s *Handlers) CreateFaq(ctx context.Context, req *pb.CreateFaqRequest) (*pb.CreateFaqResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserOnly)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &faq.CreateFaqHandler{Model: s.Model}
	resp, err := handler.CreateFaq(ctx, req)
	if err != nil {
		logger.Log.Error("CreateFaqHandler: "+err.Error(), zap.String("UserID", u.ID))
		return nil, err
	}
	logger.Log.Info("CreateFaqHandler", zap.String("UserID", u.ID))
	return resp, nil
}

func (s *Handlers) GetFaqs(ctx context.Context, req *pb.GetFaqsRequest) (*pb.GetFaqsResponse, error) {
	handler := &faq.GetFaqsHandler{Model: s.Model}
	resp, err := handler.GetFaqs(ctx, req)
	if err != nil {
		logger.Log.Error("GetFaqsHandler: " + err.Error())
		return nil, err
	}
	return resp, nil
}

func (s *Handlers) GetFaq(ctx context.Context, req *pb.GetFaqRequest) (*pb.GetFaqResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserOnly)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &faq.GetFaqHandler{Model: s.Model}
	resp, err := handler.GetFaq(ctx, req)
	if err != nil {
		logger.Log.Error("GetFaqHandler: "+err.Error(), zap.String("UserID", u.ID), zap.String("FaqID", req.Id))
		return nil, err
	}
	logger.Log.Info("GetFaqHandler", zap.String("UserID", u.ID), zap.String("FaqID", req.Id))
	return resp, nil
}

func (s *Handlers) DeleteFaq(ctx context.Context, req *pb.DeleteFaqRequest) (*pb.DeleteFaqResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserOnly)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &faq.DeleteFaqHandler{Model: s.Model}
	resp, err := handler.DeleteFaq(ctx, req)
	if err != nil {
		logger.Log.Error("DeleteFaqHandler: "+err.Error(), zap.String("UserID", u.ID), zap.String("TargetFaqID", req.Id))
		return nil, err
	}
	logger.Log.Info("DeleteFaqHandler", zap.String("UserID", u.ID), zap.String("TargetFaqID", req.Id))
	return resp, nil
}

func (s *Handlers) UpdateFaq(ctx context.Context, req *pb.UpdateFaqRequest) (*pb.UpdateFaqResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserOnly)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &faq.UpdateFaqHandler{Model: s.Model}
	resp, err := handler.UpdateFaq(ctx, req)
	if err != nil {
		logger.Log.Error("UpdateFaqHandler: "+err.Error(), zap.String("UserID", u.ID), zap.String("TargetFaqID", req.Id))
		return nil, err
	}
	logger.Log.Info("UpdateFaqHandler", zap.String("UserID", u.ID), zap.String("TargetFaqID", req.Id))
	return resp, nil
}

func (s *Handlers) DeleteFaqs(ctx context.Context, req *pb.DeleteFaqsRequest) (*pb.DeleteFaqsResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserOnly)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &faq.DeleteFaqsHandler{Model: s.Model}
	resp, err := handler.DeleteFaqs(ctx, req)
	if err != nil {
		logger.Log.Error("DeleteFaqsHandler: "+err.Error(), zap.String("UserID", u.ID), zap.Strings("TargetFaqIDs", req.Ids))
		return nil, err
	}
	logger.Log.Info("DeleteFaqsHandler", zap.String("UserID", u.ID), zap.Strings("TargetFaqIDs", req.Ids))
	return resp, nil
}

func (s *Handlers) UpdateFaqs(ctx context.Context, req *pb.UpdateFaqsRequest) (*pb.UpdateFaqsResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserOnly)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &faq.UpdateFaqsHandler{Model: s.Model}
	resp, err := handler.UpdateFaqs(ctx, req)
	if err != nil {
		logger.Log.Error("UpdateFaqsHandler: "+err.Error(), zap.String("UserID", u.ID), zap.Strings("TargetFaqIDs", req.Ids))
		return nil, err
	}
	logger.Log.Info("UpdateFaqsHandler", zap.String("UserID", u.ID), zap.Strings("TargetFaqIDs", req.Ids))
	return resp, nil
}

// -------------------- Faqs ------------------------

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

// -------------------- Activities ------------------------

func (s *Handlers) GetActivities(ctx context.Context, req *pb.GetActivitiesRequest) (*pb.GetActivitiesResponse, error) {
	u, err := s.validateUser(ctx, constants.SuperUserAndAdmin)
	if err != nil {
		return nil, constants.UnauthorizedAccessError
	}
	handler := &activity.GetActivitiesHandler{Model: s.Model}
	resp, err := handler.GetActivities(ctx, req)
	if err != nil {
		logger.Log.Error("GetActivitiesHandler: "+err.Error(), zap.String("UserID", u.ID))
		return nil, err
	}
	logger.Log.Info("GetActivitiesHandler", zap.String("UserID", u.ID))
	return resp, nil
}

// -------------------- Activities ------------------------
