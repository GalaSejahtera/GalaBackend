package user

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/logger"
	"galasejahtera/pkg/model"

	"github.com/golang/protobuf/ptypes/empty"
)

type UpdatePasswordHandler struct {
	Model model.IModel
}

func (s *UpdatePasswordHandler) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*empty.Empty, error) {
	// get user
	user, err := s.Model.GetUser(ctx, req.UserId)
	if err != nil {
		logger.Log.Error("UpdatePassword: " + err.Error())
		return nil, constants.UserNotFoundError
	}
	user.Password = req.Password

	// update user password
	_, err = s.Model.UpdateUserPassword(ctx, user)
	if err != nil {
		logger.Log.Error("UpdatePassword: " + err.Error())
		return nil, constants.InternalError
	}

	return &empty.Empty{}, nil
}
