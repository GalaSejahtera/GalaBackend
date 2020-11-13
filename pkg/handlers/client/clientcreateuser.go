package client

import (
	"context"
	pb "safeworkout/pkg/api"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/dto"
	"safeworkout/pkg/logger"
	"safeworkout/pkg/model"
	"safeworkout/pkg/utility"

	"github.com/twinj/uuid"
)

type CreateUserClientHandler struct {
	Model model.IModel
}

func (s *CreateUserClientHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if req.Data == nil {
		return nil, constants.InvalidArgumentError
	}
	user := &dto.User{
		ID:          uuid.NewV4().String(),
		Role:        constants.User,
		Name:        utility.RemoveZeroWidth(req.Data.Name),
		IC:          utility.RemoveZeroWidth(req.Data.Ic),
		PhoneNumber: utility.RemoveZeroWidth(req.Data.PhoneNumber),
		Email:       utility.RemoveZeroWidth(req.Data.Email),
		Password:    utility.RemoveZeroWidth(req.Data.Password),
		Alert:       true,
		Consent:     req.Data.Consent,
		Users:       []*dto.User{},
		Zones:       []*dto.Zone{},
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

func (s *CreateUserClientHandler) validateAndProcessReq(user *dto.User) error {
	user.Name = utility.NormalizeName(user.Name)
	user.PhoneNumber = utility.NormalizePhoneNumber(user.PhoneNumber, "")
	user.IC = utility.NormalizeID(user.IC)
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
	if user.Consent == 0 {
		return constants.ConsentNotSignedError
	}
	return nil
}

func (s *CreateUserClientHandler) userToResp(user *dto.User) *pb.CreateUserResponse {
	return &pb.CreateUserResponse{
		Data: &pb.User{
			Id:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			IsActive:    user.IsActive,
			LastUpdated: user.LastUpdated,
			Lat:         user.Lat,
			Long:        user.Long,
			Ic:          user.IC,
			PhoneNumber: user.PhoneNumber,
			Alert:       user.Alert,
		},
	}
}
