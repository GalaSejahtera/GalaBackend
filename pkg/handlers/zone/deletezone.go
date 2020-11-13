package zone

import (
	"context"
	pb "safeworkout/pkg/api"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/dto"
	"safeworkout/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteZoneHandler struct {
	Model model.IModel
}

func (s *DeleteZoneHandler) DeleteZone(ctx context.Context, req *pb.DeleteZoneRequest) (*pb.DeleteZoneResponse, error) {
	rslt, err := s.Model.DeleteZone(ctx, req.Id)
	if err != nil {
		if status.Code(err) == codes.Unknown {
			return nil, constants.ZoneNotFoundError
		}
		return nil, constants.InternalError
	}

	resp := s.zoneToResp(rslt)
	return resp, nil
}

func (s *DeleteZoneHandler) zoneToResp(zone *dto.Zone) *pb.DeleteZoneResponse {
	resp := &pb.DeleteZoneResponse{
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
