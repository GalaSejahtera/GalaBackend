package zone

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/model"
	"galasejahtera/pkg/utility"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateZoneHandler struct {
	Model model.IModel
}

func (s *UpdateZoneHandler) UpdateZone(ctx context.Context, req *pb.UpdateZoneRequest) (*pb.UpdateZoneResponse, error) {
	if req.Data == nil {
		return nil, constants.InvalidArgumentError
	}
	zone := s.reqToZone(req)

	_, err := s.Model.GetZone(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	v, err := s.Model.UpdateZone(ctx, zone)
	if err != nil {
		if status.Code(err) == codes.Unknown {
			return nil, constants.ZoneNotFoundError
		}
		return nil, constants.InternalError
	}
	resp := s.zoneToResp(v)
	return resp, nil
}

func (s *UpdateZoneHandler) reqToZone(req *pb.UpdateZoneRequest) *dto.Zone {
	zone := &dto.Zone{
		ID:       req.Id,
		Name:     utility.RemoveZeroWidth(req.Data.Name),
		Lat:      req.Data.Lat,
		Long:     req.Data.Long,
		Type:     req.Data.Type,
		Capacity: req.Data.Capacity,
		Radius:   req.Data.Radius,
	}
	return zone
}

func (s *UpdateZoneHandler) zoneToResp(zone *dto.Zone) *pb.UpdateZoneResponse {
	resp := &pb.UpdateZoneResponse{
		Data: &pb.Zone{
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
		},
	}
	return resp
}
