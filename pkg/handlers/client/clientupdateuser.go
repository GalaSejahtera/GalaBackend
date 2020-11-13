package client

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/logger"
	"galasejahtera/pkg/model"
	"galasejahtera/pkg/utility"
)

type UpdateUserClientHandler struct {
	Model model.IModel
}

func (s *UpdateUserClientHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if req.Data == nil {
		return nil, constants.InvalidArgumentError
	}
	user := s.reqToUser(req)

	err := s.validateAndProcessReq(user)
	if err != nil {
		return nil, err
	}

	_, err = s.Model.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	v, err := s.Model.ClientUpdateAppUser(ctx, user)
	if err != nil {
		logger.Log.Error("UpdateUserClient: " + err.Error())
		return nil, constants.InternalError
	}
	resp := s.userToResp(v)
	return resp, nil
}

func (s *UpdateUserClientHandler) reqToUser(req *pb.UpdateUserRequest) *dto.User {
	user := &dto.User{
		ID:          utility.RemoveZeroWidth(req.Id),
		Name:        utility.RemoveZeroWidth(req.Data.Name),
		IC:          utility.RemoveZeroWidth(req.Data.Ic),
		PhoneNumber: utility.RemoveZeroWidth(req.Data.PhoneNumber),
		Alert:       req.Data.Alert,
		Infected:    req.Data.Infected,
	}
	return user
}

func (s *UpdateUserClientHandler) userToResp(user *dto.User) *pb.UpdateUserResponse {
	resp := &pb.UpdateUserResponse{
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
	return resp
}

func (s *UpdateUserClientHandler) validateAndProcessReq(user *dto.User) error {
	user.Name = utility.NormalizeName(user.Name)
	oldPhoneNumber := user.PhoneNumber
	user.PhoneNumber = utility.NormalizePhoneNumber(user.PhoneNumber, "")
	user.IC = utility.NormalizeID(user.IC)
	if oldPhoneNumber != "" && user.PhoneNumber == "" {
		return constants.InvalidPhoneNumberError
	}
	return nil
}
