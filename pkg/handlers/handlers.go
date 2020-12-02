package handlers

import (
	"context"
	"errors"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/handlers/faq"
	"galasejahtera/pkg/handlers/kase"
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

// -------------------- Faqs ------------------------

func (s *Handlers) CreateFaq(ctx context.Context, req *pb.CreateFaqRequest) (*pb.CreateFaqResponse, error) {
	handler := &faq.CreateFaqHandler{Model: s.Model}
	resp, err := handler.CreateFaq(ctx, req)
	if err != nil {
		logger.Log.Error("CreateFaqHandler: " + err.Error())
		return nil, err
	}
	logger.Log.Info("CreateFaqHandler")
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
	u, err := s.validateUser(ctx, constants.AllCanAccess)
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
	u, err := s.validateUser(ctx, constants.AllCanAccess)
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
	u, err := s.validateUser(ctx, constants.AllCanAccess)
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
	u, err := s.validateUser(ctx, constants.AllCanAccess)
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
	u, err := s.validateUser(ctx, constants.AllCanAccess)
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
