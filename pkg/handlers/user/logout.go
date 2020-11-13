package user

import (
	"context"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/logger"
	"galasejahtera/pkg/model"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/metadata"
)

type LogoutHandler struct {
	Model model.IModel
}

func (s *LogoutHandler) Logout(ctx context.Context) (*empty.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, constants.MetadataNotFoundError
	}
	tokenSlice := md.Get("authorization")
	err := s.Model.Logout(ctx, strings.Join(tokenSlice, " "))
	if err != nil {
		logger.Log.Error("Logout: " + err.Error())
	}
	return &empty.Empty{}, nil
}
