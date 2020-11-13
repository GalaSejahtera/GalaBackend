package zone

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/logger"
	"galasejahtera/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetZoneHandler struct {
	Model model.IModel
}

func (s *GetZoneHandler) GetZone(ctx context.Context, req *pb.GetZoneRequest) (*pb.GetZoneResponse, error) {
	zone, err := s.Model.GetZone(ctx, req.Id)
	if err != nil {
		logger.Log.Error("GetZoneHandler: " + err.Error())
		if status.Code(err) == codes.Unknown {
			return nil, constants.ZoneNotFoundError
		}
		return nil, constants.InternalError
	}

	resp, err := s.zoneToResponse(zone)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GetZoneHandler) zoneToResponse(zone *dto.Zone) (*pb.GetZoneResponse, error) {
	resp := &pb.GetZoneResponse{
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

	return resp, nil
}
