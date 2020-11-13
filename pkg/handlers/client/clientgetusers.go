package client

import (
	"context"
	pb "safeworkout/pkg/api"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/dto"
	"safeworkout/pkg/logger"
	"safeworkout/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetUsersClientHandler struct {
	Model model.IModel
}

func (s *GetUsersClientHandler) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	var sort *dto.SortData
	var itemsRange *dto.RangeData
	var filter *dto.FilterData

	// If the request is batch get, call batch get model
	if len(req.Ids) > 0 {
		users, err := s.Model.BatchGetUsers(ctx, req.Ids)
		if err != nil {
			if status.Code(err) == codes.Unknown {
				return nil, constants.UserNotFoundError
			}
			return nil, constants.InternalError
		}
		resp, err := s.usersToResponses(users)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}

	if req.Item != "" && req.Order != "" {
		sort = &dto.SortData{
			Item:  req.Item,
			Order: req.Order,
		}
	}

	if req.To != 0 {
		itemsRange = &dto.RangeData{
			From: int(req.From),
			To:   int(req.To),
		}
	}

	if req.FilterItem != "" && req.FilterValue != "" {
		filter = &dto.FilterData{
			Item:  req.FilterItem,
			Value: req.FilterValue,
		}
	}

	total, users, err := s.Model.QueryUsers(ctx, sort, itemsRange, filter)
	if err != nil {
		logger.Log.Error("GetUsersClientHandler: " + err.Error())
		if status.Code(err) == codes.Unknown {
			return nil, constants.UserNotFoundError
		}
		return nil, constants.InternalError
	}

	resp, err := s.usersToResponses(users)
	if err != nil {
		return nil, err
	}

	resp.Total = total
	return resp, nil
}

func (s *GetUsersClientHandler) usersToResponses(users []*dto.User) (*pb.GetUsersResponse, error) {
	var resps []*pb.User
	for _, user := range users {
		resp := &pb.User{
			Id:   user.ID,
			Lat:  user.Lat,
			Long: user.Long,
		}

		resps = append(resps, resp)
	}
	rslt := &pb.GetUsersResponse{
		Data: resps,
	}

	return rslt, nil
}
