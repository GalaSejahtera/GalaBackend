package user

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/logger"
	"galasejahtera/pkg/model"
)

type LoginHandler struct {
	Model model.IModel
}

func (s *LoginHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.Model.Login(ctx, req.Email, req.Password)
	if err != nil {
		logger.Log.Error("Login: " + err.Error())
		return nil, err
	}
	resp := s.userToResp(user)
	return resp, nil
}

func (s *LoginHandler) userToResp(user *dto.User) *pb.LoginResponse {
	return &pb.LoginResponse{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		Role:         user.Role,
		Id:           user.ID,
	}
}
