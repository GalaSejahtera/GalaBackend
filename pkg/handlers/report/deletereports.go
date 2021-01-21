package report

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteReportsHandler struct {
	Model model.IModel
}

func (s *DeleteReportsHandler) DeleteReports(ctx context.Context, req *pb.DeleteReportsRequest) (*pb.DeleteReportsResponse, error) {
	var ids []string

	// remove reports
	for _, id := range req.Ids {
		u, err := s.Model.DeleteReport(ctx, id)
		if err != nil {
			if status.Code(err) == codes.Unknown {
				return nil, constants.ReportNotFoundError
			}
			return nil, constants.InternalError
		}

		// add report into deleted report IDs
		ids = append(ids, u.ID)
	}

	return &pb.DeleteReportsResponse{Data: ids}, nil
}
