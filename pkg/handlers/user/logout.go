package user

import (
	"context"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/logger"
	"safeworkout/pkg/model"
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
