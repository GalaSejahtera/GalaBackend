package client

import (
	"context"
	pb "safeworkout/pkg/api"
	"safeworkout/pkg/dto"
	"safeworkout/pkg/logger"
	"safeworkout/pkg/model"
)

type GetRecentZonesHandler struct {
	Model model.IModel
}

func (s *GetRecentZonesHandler) GetRecentZones(ctx context.Context, req *pb.ClientGetRecentZonesRequest) (*pb.GetZonesResponse, error) {
	zones, err := s.Model.GetRecentZonesByUser(ctx, req.Id)
	if err != nil {
		logger.Log.Error("GetRecentZonesHandler: " + err.Error())
		return nil, err
	}

	resp, err := s.zonesToResponses(zones)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GetRecentZonesHandler) zonesToResponses(zones []*dto.Zone) (*pb.GetZonesResponse, error) {
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
			Time:               zone.Time,
			Risk:               zone.Risk,
		}

		resps = append(resps, resp)
	}
	rslt := &pb.GetZonesResponse{
		Data: resps,
	}

	return rslt, nil
}
