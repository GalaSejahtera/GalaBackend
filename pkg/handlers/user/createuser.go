package user

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/logger"
	"galasejahtera/pkg/model"
	"galasejahtera/pkg/utility"

	"github.com/twinj/uuid"
)

type CreateUserHandler struct {
	Model model.IModel
}

func (s *CreateUserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if req.Data == nil {
		return nil, constants.InvalidArgumentError
	}
	user := &dto.User{
		ID:       uuid.NewV4().String(),
		Role:     utility.RemoveZeroWidth(req.Data.Role),
		Email:    utility.RemoveZeroWidth(req.Data.Email),
		Password: utility.RemoveZeroWidth(req.Data.Password),
		Lat:      req.Data.Lat,
		Long:     req.Data.Long,
		Users:    []*dto.User{},
		IsActive: req.Data.IsActive,
		Name:     req.Data.Name,
	}
	err := s.validateAndProcessReq(user)
	if err != nil {
		return nil, err
	}

	// check if email exist
	count, _, err := s.Model.QueryUsers(ctx, nil, nil, &dto.FilterData{
		Item:  constants.Email,
		Value: user.Email,
	})
	if count > 0 {
		return nil, constants.EmailAlreadyExistError
	}

	rslt, err := s.Model.CreateUser(ctx, user)
	if err != nil {
		logger.Log.Error("CreateUser: " + err.Error())
		return nil, constants.InternalError
	}
	resp := s.userToResp(rslt)
	return resp, nil
}

func (s *CreateUserHandler) validateAndProcessReq(user *dto.User) error {
	valid := utility.ValidateEmail(user.Email)
	if !valid {
		return constants.InvalidEmailError
	}
	if len(user.Password) < 6 {
		return constants.InvalidPasswordError
	}
	if user.Email == "" {
		return constants.InvalidEmailError
	}
	if user.Role == "" {
		return constants.InvalidRoleError
	}

	return nil
}

func (s *CreateUserHandler) userToResp(user *dto.User) *pb.CreateUserResponse {
	return &pb.CreateUserResponse{
		Data: &pb.User{
			Id:          user.ID,
			Role:        user.Role,
			Email:       user.Email,
			LastUpdated: user.LastUpdated,
			Lat:         user.Lat,
			Long:        user.Long,
			IsActive:    user.IsActive,
			Name:        user.Name,
		},
	}
}
