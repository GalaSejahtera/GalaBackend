package user

import (
	"context"
	pb "safeworkout/pkg/api"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteUsersHandler struct {
	Model model.IModel
}

func (s *DeleteUsersHandler) DeleteUsers(ctx context.Context, req *pb.DeleteUsersRequest) (*pb.DeleteUsersResponse, error) {
	var ids []string

	// remove users
	for _, id := range req.Ids {
		u, err := s.Model.DeleteUser(ctx, id)
		if err != nil {
			if status.Code(err) == codes.Unknown {
				return nil, constants.UserNotFoundError
			}
			return nil, constants.InternalError
		}

		// add user into deleted user IDs
		ids = append(ids, u.ID)
	}

	return &pb.DeleteUsersResponse{Data: ids}, nil
}
