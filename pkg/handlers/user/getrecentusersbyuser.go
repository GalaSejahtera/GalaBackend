package user

import (
	"context"
	pb "safeworkout/pkg/api"
	"safeworkout/pkg/dto"
	"safeworkout/pkg/model"
	"sort"
)

type GetRecentUsersByUserHandler struct {
	Model model.IModel
}

func (s *GetRecentUsersByUserHandler) GetRecentUsersByUser(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {

	users, err := s.Model.GetRecentUsersByUser(ctx, req.FilterValue)
	if err != nil {
		return nil, err
	}

	if req.Item != "" && req.Order != "" {
		s.sortUsers(users, req.Item, req.Order)
	}

	// get total
	total := int64(len(users))

	index := int64(0)
	var rslt []*dto.User

	for _, user := range users {
		if index < req.From {
			index += 1
			continue
		}
		if req.To != 0 && index > req.To {
			index += 1
			continue
		}
		index += 1
		rslt = append(rslt, user)
	}

	resp := s.usersToResponses(rslt)
	resp.Total = total

	return resp, nil
}

func (s *GetRecentUsersByUserHandler) usersToResponses(users []*dto.User) *pb.GetUsersResponse {
	var resps []*pb.User
	for _, user := range users {
		resp := &pb.User{
			Id:          user.ID,
			Role:        user.Role,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
			Ic:          user.IC,
			Email:       user.Email,
			IsActive:    user.IsActive,
			LastUpdated: user.LastUpdated,
			Lat:         user.Lat,
			Long:        user.Long,
			Consent:     user.Consent,
			Time:        user.Time,
			Infected:    user.Infected,
		}

		resps = append(resps, resp)
	}
	rslt := &pb.GetUsersResponse{
		Data: resps,
	}

	return rslt
}

func (s *GetRecentUsersByUserHandler) sortUsers(users []*dto.User, field string, order string) {
	switch field {
	case "id":
		sort.Slice(users, func(i, j int) bool {
			return users[i].ID < users[j].ID
		})
	case "name":
		sort.Slice(users, func(i, j int) bool {
			return users[i].Name < users[j].Name
		})
	case "ic":
		sort.Slice(users, func(i, j int) bool {
			return users[i].IC < users[j].IC
		})
	case "phoneNumber":
		sort.Slice(users, func(i, j int) bool {
			return users[i].PhoneNumber < users[j].PhoneNumber
		})
	case "email":
		sort.Slice(users, func(i, j int) bool {
			return users[i].Email < users[j].Email
		})
	case "lat":
		sort.Slice(users, func(i, j int) bool {
			return users[i].Lat < users[j].Lat
		})
	case "long":
		sort.Slice(users, func(i, j int) bool {
			return users[i].Long < users[j].Long
		})
	case "lastUpdated":
		sort.Slice(users, func(i, j int) bool {
			return users[i].LastUpdated < users[j].LastUpdated
		})
	case "time":
		sort.Slice(users, func(i, j int) bool {
			return users[i].Time < users[j].Time
		})
	default:
	}

	if order == "DESC" {
		// reverse slice
		for i, j := 0, len(users)-1; i < j; i, j = i+1, j-1 {
			users[i], users[j] = users[j], users[i]
		}
	}
}
