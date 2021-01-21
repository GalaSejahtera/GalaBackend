package user

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/logger"
	"galasejahtera/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetUserHandler struct {
	Model model.IModel
}

func (s *GetUserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.Model.GetUser(ctx, req.Id)
	if err != nil {
		logger.Log.Error("GetUserHandler: " + err.Error())
		if status.Code(err) == codes.Unknown {
			return nil, constants.UserNotFoundError
		}
		return nil, constants.InternalError
	}

	resp, err := s.userToResponse(user)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GetUserHandler) userToResponse(user *dto.User) (*pb.GetUserResponse, error) {
	resp := &pb.GetUserResponse{
		Data: &pb.User{
			Id:          user.ID,
			Role:        user.Role,
			Email:       user.Email,
			IsActive:    user.IsActive,
			LastUpdated: user.LastUpdated,
			Lat:         user.Lat,
			Long:        user.Long,
			Name:        user.Name,
		},
	}

	return resp, nil
}
