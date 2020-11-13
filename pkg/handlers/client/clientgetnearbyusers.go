package client

import (
	"context"
	pb "safeworkout/pkg/api"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/dto"
	"safeworkout/pkg/logger"
	"safeworkout/pkg/model"
	"safeworkout/pkg/utility"
	"time"
)

type GetNearbyUsersHandler struct {
	Model model.IModel
}

func (s *GetNearbyUsersHandler) GetNearbyUsers(ctx context.Context, req *pb.ClientGetNearbyUsersRequest) (*pb.ClientGetNearbyUsersResponse, error) {
	if req.User == nil {
		return nil, constants.InvalidArgumentError
	}

	u := &dto.User{
		ID:          req.User.Id,
		IsActive:    true,
		LastUpdated: utility.TimeToMilli(utility.MalaysiaTime(time.Now())),
		Lat:         req.User.Lat,
		Long:        req.User.Long,
		Location: &dto.Location{
			Type:        "Point",
			Coordinates: []float64{req.User.Long, req.User.Lat},
		},
	}

	// get users
	total, users, err := s.Model.GetNearbyUsers(ctx, u)
	if err != nil {
		logger.Log.Error("GetNearbyUsersHandler: " + err.Error())
		return nil, err
	}

	resp, err := s.usersToResponses(users)
	if err != nil {
		return nil, err
	}

	resp.UserNum = total

	return resp, nil
}

func (s *GetNearbyUsersHandler) usersToResponses(users []*dto.User) (*pb.ClientGetNearbyUsersResponse, error) {
	var resps []*pb.User
	for _, user := range users {
		resp := &pb.User{
			Lat:  user.Lat,
			Long: user.Long,
		}

		resps = append(resps, resp)
	}
	rslt := &pb.ClientGetNearbyUsersResponse{
		Users: resps,
	}

	return rslt, nil
}
