package zone

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteZonesHandler struct {
	Model model.IModel
}

func (s *DeleteZonesHandler) DeleteZones(ctx context.Context, req *pb.DeleteZonesRequest) (*pb.DeleteZonesResponse, error) {
	var ids []string

	// remove zones
	for _, id := range req.Ids {
		u, err := s.Model.DeleteZone(ctx, id)
		if err != nil {
			if status.Code(err) == codes.Unknown {
				return nil, constants.ZoneNotFoundError
			}
			return nil, constants.InternalError
		}

		// add zone into deleted zone IDs
		ids = append(ids, u.ID)
	}

	return &pb.DeleteZonesResponse{Data: ids}, nil
}
