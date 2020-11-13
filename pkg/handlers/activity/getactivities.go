package activity

import (
	"context"
	pb "safeworkout/pkg/api"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/dto"
	"safeworkout/pkg/logger"
	"safeworkout/pkg/model"
	"safeworkout/pkg/utility"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetActivitiesHandler struct {
	Model model.IModel
}

func (s *GetActivitiesHandler) GetActivities(ctx context.Context, req *pb.GetActivitiesRequest) (*pb.GetActivitiesResponse, error) {
	var sort *dto.SortData
	var itemsRange *dto.RangeData
	var filter *dto.FilterData

	// If the request is get by id, call get model
	if len(req.Ids) > 0 {
		activity, err := s.Model.GetActivity(ctx, req.Ids[0])
		if err != nil {
			if status.Code(err) == codes.Unknown {
				return nil, constants.ActivityNotFoundError
			}
			return nil, constants.InternalError
		}
		resp, err := s.activitiesToResponses([]*dto.Activity{activity})
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

	total, activities, err := s.Model.QueryActivities(ctx, sort, itemsRange, filter)
	if err != nil {
		logger.Log.Error("GetActivitiesHandler: " + err.Error())
		if status.Code(err) == codes.Unknown {
			return nil, constants.ActivityNotFoundError
		}
		return nil, constants.InternalError
	}

	resp, err := s.activitiesToResponses(activities)
	if err != nil {
		return nil, err
	}

	resp.Total = total
	return resp, nil
}

func (s *GetActivitiesHandler) activitiesToResponses(activities []*dto.Activity) (*pb.GetActivitiesResponse, error) {
	var resps []*pb.Activity
	for _, activity := range activities {
		resp := &pb.Activity{
			Id:     activity.ID,
			ZoneId: activity.ZoneID,
			UserId: activity.UserID,
			Time:   activity.Time,
			Ttl:    utility.TimeToMilli(activity.TTL),
		}

		resps = append(resps, resp)
	}
	rslt := &pb.GetActivitiesResponse{
		Data: resps,
	}

	return rslt, nil
}
