package user

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/logger"
	"galasejahtera/pkg/model"
	"galasejahtera/pkg/utility"
	"time"
)

type GetNearbyUsersHandler struct {
	Model model.IModel
}

func (s *GetNearbyUsersHandler) GetNearbyUsers(ctx context.Context, req *pb.GetNearbyUsersRequest) (*pb.GetNearbyUsersResponse, error) {
	if req.User == nil {
		return nil, constants.InvalidArgumentError
	}

	u, err := s.Model.GetUser(ctx, req.User.Id)
	if err != nil {
		return nil, err
	}

	// patch user
	u.IsActive = true
	u.LastUpdated = utility.TimeToMilli(utility.MalaysiaTime(time.Now()))
	u.Lat = req.User.Lat
	u.Long = req.User.Long
	u.Location = &dto.Location{
		Type:        "Point",
		Coordinates: []float64{req.User.Long, req.User.Lat},
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

func (s *GetNearbyUsersHandler) usersToResponses(users []*dto.User) (*pb.GetNearbyUsersResponse, error) {
	var resps []*pb.User
	for _, user := range users {
		resp := &pb.User{
			Lat:  user.Lat,
			Long: user.Long,
		}

		resps = append(resps, resp)
	}
	rslt := &pb.GetNearbyUsersResponse{
		Users: resps,
	}

	return rslt, nil
}
