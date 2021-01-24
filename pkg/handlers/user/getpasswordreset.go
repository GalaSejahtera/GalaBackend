package user

import (
	"context"
	"fmt"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/logger"
	"galasejahtera/pkg/model"
	"galasejahtera/pkg/utility"
	"github.com/twinj/uuid"
)

type GetPasswordResetHandler struct {
	Model model.IModel
}

func (s *GetPasswordResetHandler) GetPasswordReset(ctx context.Context, req *pb.GetPasswordResetRequest) (*pb.GetPasswordResetResponse, error) {
	// get user
	_, us, err := s.Model.QueryUsers(ctx, nil, nil, &dto.FilterData{
		Item:  constants.Email,
		Value: req.Email,
	})
	if err != nil {
		return nil, err
	}
	if len(us) < 1 {
		return nil, constants.UserNotFoundError
	}
	user := us[0]

	// generate random password
	user.Password = utility.NormalizeContent(uuid.NewV4().String())

	// update user password
	_, err = s.Model.UpdateUserPassword(ctx, user)
	if err != nil {
		return nil, err
	}

	// send password reset email
	err = utility.SendPasswordResetEmail(user.Email, user.Email, user.Password)
	if err != nil {
		logger.Log.Warn("GetPasswordReset: " + err.Error())
		return nil, err
	}

	return &pb.GetPasswordResetResponse{
		Message: fmt.Sprintf("New password has been sent to the email %s, please check your inbox and spam box.", user.Email),
	}, nil
}
