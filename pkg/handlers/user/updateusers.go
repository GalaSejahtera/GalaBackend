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

type UpdateUsersHandler struct {
	Model model.IModel
}

func (s *UpdateUsersHandler) UpdateUsers(ctx context.Context, req *pb.UpdateUsersRequest) (*pb.UpdateUsersResponse, error) {
	if req.Data == nil {
		return nil, constants.InvalidArgumentError
	}
	if len(req.Ids) > 1 {
		return nil, constants.OperationUnsupportedError
	}
	user := s.reqToUser(req)
	err := s.validateAndProcessReq(user)
	if err != nil {
		return nil, err
	}

	// get user
	u, err := s.Model.GetUser(ctx, req.Ids[0])
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

	ids, err := s.Model.UpdateUsers(ctx, user, req.Ids)
	if err != nil {
		logger.Log.Error("UpdateUsers: " + err.Error())
		return nil, constants.InternalError
	}
	return &pb.UpdateUsersResponse{Data: ids}, nil
}

func (s *UpdateUsersHandler) reqToUser(req *pb.UpdateUsersRequest) *dto.User {
	user := &dto.User{
		Role:        utility.RemoveZeroWidth(req.Data.Role),
		Name:        utility.RemoveZeroWidth(req.Data.Name),
		IC:          utility.RemoveZeroWidth(req.Data.Ic),
		PhoneNumber: utility.RemoveZeroWidth(req.Data.PhoneNumber),
		Email:       utility.RemoveZeroWidth(req.Data.Email),
		Password:    utility.RemoveZeroWidth(req.Data.Password),
		IsActive:    req.Data.IsActive,
		Lat:         req.Data.Lat,
		Long:        req.Data.Long,
	}
	return user
}

func (s *UpdateUsersHandler) validateAndProcessReq(user *dto.User) error {
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
