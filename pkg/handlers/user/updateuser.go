package user

import (
	"context"
	pb "safeworkout/pkg/api"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/dto"
	"safeworkout/pkg/logger"
	"safeworkout/pkg/model"
	"safeworkout/pkg/utility"
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
		ID:          utility.RemoveZeroWidth(req.Id),
		Role:        utility.RemoveZeroWidth(req.Data.Role),
		Name:        utility.RemoveZeroWidth(req.Data.Name),
		IC:          utility.RemoveZeroWidth(req.Data.Ic),
		PhoneNumber: utility.RemoveZeroWidth(req.Data.PhoneNumber),
		Email:       utility.RemoveZeroWidth(req.Data.Email),
		IsActive:    req.Data.IsActive,
		Lat:         req.Data.Lat,
		Long:        req.Data.Long,
		Infected:    req.Data.Infected,
	}
	return user
}

func (s *UpdateUserHandler) userToResp(user *dto.User) *pb.UpdateUserResponse {
	resp := &pb.UpdateUserResponse{
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

func (s *UpdateUserHandler) validateAndProcessReq(user *dto.User) error {
	user.Name = utility.NormalizeName(user.Name)
	user.PhoneNumber = utility.NormalizePhoneNumber(user.PhoneNumber, "")
	user.IC = utility.NormalizeID(user.IC)
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
