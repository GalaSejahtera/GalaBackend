package zone

import (
	"context"
	pb "safeworkout/pkg/api"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/dto"
	"safeworkout/pkg/model"
	"safeworkout/pkg/utility"

	"github.com/twinj/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateZoneHandler struct {
	Model model.IModel
}

func (s *CreateZoneHandler) CreateZone(ctx context.Context, req *pb.CreateZoneRequest) (*pb.CreateZoneResponse, error) {
	if req.Data == nil {
		return nil, constants.InvalidArgumentError
	}
	zone := &dto.Zone{
		ID:       uuid.NewV4().String(),
		Name:     utility.RemoveZeroWidth(req.Data.Name),
		Lat:      req.Data.Lat,
		Long:     req.Data.Long,
		Type:     req.Data.Type,
		Capacity: req.Data.Capacity,
		Radius:   req.Data.Radius,
	}

	rslt, err := s.Model.CreateZone(ctx, zone)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return nil, constants.ZoneAlreadyExistError
		}
		return nil, constants.InternalError
	}
	resp := s.zoneToResp(rslt)
	return resp, nil
}

func (s *CreateZoneHandler) zoneToResp(zone *dto.Zone) *pb.CreateZoneResponse {
	return &pb.CreateZoneResponse{
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
}
