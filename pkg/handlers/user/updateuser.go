package user

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/logger"
	"galasejahtera/pkg/model"
	"galasejahtera/pkg/utility"
)

type UpdateUserHandler struct {
	Model model.IModel
}

func (s *UpdateUserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if req.Data == nil {
		return nil, constants.InvalidArgumentError
	}
	user := s.reqToUser(req)

	err := s.validateAndProcessReq(user)
	if err != nil {
		return nil, err
	}

	u, err := s.Model.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if u.Email != user.Email {
		// check if email exist
		count, _, err := s.Model.QueryUsers(ctx, nil, nil, &dto.FilterData{
			Item:  constants.Email,
			Value: user.Email,
		})
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, constants.EmailAlreadyExistError
		}
	}

	v, err := s.Model.UpdateUser(ctx, user)
	if err != nil {
		logger.Log.Error("UpdateUser: " + err.Error())
		return nil, constants.InternalError
	}
	resp := s.userToResp(v)
	return resp, nil
}

func (s *UpdateUserHandler) reqToUser(req *pb.UpdateUserRequest) *dto.User {
	user := &dto.User{
		ID:       utility.RemoveZeroWidth(req.Id),
		Role:     utility.RemoveZeroWidth(req.Data.Role),
		Email:    utility.RemoveZeroWidth(req.Data.Email),
		IsActive: req.Data.IsActive,
		Lat:      req.Data.Lat,
		Long:     req.Data.Long,
		Name:     req.Data.Name,
	}
	return user
}

func (s *UpdateUserHandler) userToResp(user *dto.User) *pb.UpdateUserResponse {
	resp := &pb.UpdateUserResponse{
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
	return resp
}

func (s *UpdateUserHandler) validateAndProcessReq(user *dto.User) error {
	valid := utility.ValidateEmail(user.Email)
	if !valid {
		return constants.InvalidEmailError
	}
	if user.Email == "" {
		return constants.InvalidEmailError
	}
	if user.Role == "" {
		return constants.InvalidRoleError
	}

	return nil
}
