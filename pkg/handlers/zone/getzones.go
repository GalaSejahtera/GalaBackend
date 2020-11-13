package zone

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

type GetZonesHandler struct {
	Model model.IModel
}

func (s *GetZonesHandler) GetZones(ctx context.Context, req *pb.GetZonesRequest) (*pb.GetZonesResponse, error) {
	var sort *dto.SortData
	var itemsRange *dto.RangeData
	var filter *dto.FilterData

	// If the request is batch get, call batch get model
	if len(req.Ids) > 0 {
		zones, err := s.Model.BatchGetZones(ctx, req.Ids)
		if err != nil {
			if status.Code(err) == codes.Unknown {
				return nil, constants.ZoneNotFoundError
			}
			return nil, constants.InternalError
		}
		resp, err := s.zonesToResponses(zones)
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

	total, zones, err := s.Model.QueryZones(ctx, sort, itemsRange, filter)
	if err != nil {
		logger.Log.Error("GetZonesHandler: " + err.Error())
		if status.Code(err) == codes.Unknown {
			return nil, constants.ZoneNotFoundError
		}
		return nil, constants.InternalError
	}

	resp, err := s.zonesToResponses(zones)
	if err != nil {
		return nil, err
	}

	resp.Total = total
	return resp, nil
}

func (s *GetZonesHandler) zonesToResponses(zones []*dto.Zone) (*pb.GetZonesResponse, error) {
	var resps []*pb.Zone
	for _, zone := range zones {
		resp := &pb.Zone{
			Id:                 zone.ID,
			Name:               zone.Name,
			Lat:                zone.Lat,
			Long:               zone.Long,
			Type:               zone.Type,
			Radius:             zone.Radius,
			Capacity:           zone.Capacity,
			UsersWithin:        zone.UsersWithin,
			IsCapacityExceeded: zone.IsCapacityExceeded,
			Risk:               zone.Risk,
		}

		resps = append(resps, resp)
	}
	rslt := &pb.GetZonesResponse{
		Data: resps,
	}

	return rslt, nil
}
