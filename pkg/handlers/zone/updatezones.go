package zone

import (
	"context"
	pb "safeworkout/pkg/api"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/dto"
	"safeworkout/pkg/model"
	"safeworkout/pkg/utility"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateZonesHandler struct {
	Model model.IModel
}

func (s *UpdateZonesHandler) UpdateZones(ctx context.Context, req *pb.UpdateZonesRequest) (*pb.UpdateZonesResponse, error) {
	if req.Data == nil {
		return nil, constants.InvalidArgumentError
	}
	if len(req.Ids) > 1 {
		return nil, constants.OperationUnsupportedError
	}
	zone := s.reqToZone(req)

	// get zone
	_, err := s.Model.GetZone(ctx, req.Ids[0])
	if err != nil {
		return nil, err
	}

	ids, err := s.Model.UpdateZones(ctx, zone, req.Ids)
	if err != nil {
		if status.Code(err) == codes.Unknown {
			return nil, constants.ZoneNotFoundError
		}
		return nil, constants.InternalError
	}
	return &pb.UpdateZonesResponse{Data: ids}, nil
}

func (s *UpdateZonesHandler) reqToZone(req *pb.UpdateZonesRequest) *dto.Zone {
	zone := &dto.Zone{
		Name:     utility.RemoveZeroWidth(req.Data.Name),
		Lat:      req.Data.Lat,
		Long:     req.Data.Long,
		Type:     req.Data.Type,
		Capacity: req.Data.Capacity,
		Radius:   req.Data.Radius,
	}
	return zone
}
