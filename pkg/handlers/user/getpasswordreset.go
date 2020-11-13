package user

import (
	"context"
	"fmt"
	"os"
	pb "safeworkout/pkg/api"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/dto"
	"safeworkout/pkg/logger"
	"safeworkout/pkg/model"
	"safeworkout/pkg/utility"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetPasswordResetHandler struct {
	Model model.IModel
}

func (s *GetPasswordResetHandler) GetPasswordReset(ctx context.Context, req *pb.GetPasswordResetRequest) (*pb.GetPasswordResetResponse, error) {

	user := &dto.User{}
	var err error

	// check if email is used
	if req.Email != "" {
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
		user = us[0]
	} else {
		// get user
		user, err = s.Model.GetUser(ctx, req.Id)
		if err != nil {
			if status.Code(err) == codes.Unknown {
				return nil, constants.UserNotFoundError
			}
			return nil, constants.InternalError
		}
	}

	auth, err := s.Model.CreateToken(ctx, &dto.AuthObject{
		UserId:      req.Id,
		TTL:         utility.MilliToTime(time.Now().Add(time.Hour*24*constants.PasswordResetTokenTTLDays).Unix()*1000 - 1000),
		DisplayName: user.Name,
		Type:        constants.Refresh,
	})
	if err != nil {
		logger.Log.Error("GetPasswordReset: " + err.Error())
		return nil, err
	}

	passwordReset := os.Getenv("ADMIN_URL") + "/#/resetpassword?token=" + auth.Token

	// send password reset email
	err = utility.SendPasswordResetEmail(user.Email, user.Name, passwordReset)
	if err != nil {
		logger.Log.Warn("GetPasswordReset: " + err.Error())
		return nil, err
	}

	return &pb.GetPasswordResetResponse{
		Message: fmt.Sprintf("Password reset link has been sent to the email %s, please check your inbox and spam box.", user.Email),
	}, nil
}
