package user

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteUserHandler struct {
	Model model.IModel
}

func (s *DeleteUserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	rslt, err := s.Model.DeleteUser(ctx, req.Id)
	if err != nil {
		if status.Code(err) == codes.Unknown {
			return nil, constants.UserNotFoundError
		}
		return nil, constants.InternalError
	}

	resp := s.userToResp(rslt)
	return resp, nil
}

func (s *DeleteUserHandler) userToResp(user *dto.User) *pb.DeleteUserResponse {
	resp := &pb.DeleteUserResponse{
		Data: &pb.User{
			Id:          user.ID,
			Role:        user.Role,
			Name:        user.Name,
			Ic:          user.IC,
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
			IsActive:    user.IsActive,
			LastUpdated: user.LastUpdated,
			Lat:         user.Lat,
			Long:        user.Long,
			Consent:     user.Consent,
			Infected:    user.Infected,
		},
	}
	return resp
}
