package report

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateReportsHandler struct {
	Model model.IModel
}

func (s *UpdateReportsHandler) UpdateReports(ctx context.Context, req *pb.UpdateReportsRequest) (*pb.UpdateReportsResponse, error) {
	if req.Data == nil {
		return nil, constants.InvalidArgumentError
	}
	if len(req.Ids) > 1 {
		return nil, constants.OperationUnsupportedError
	}
	report := s.reqToReport(req)

	// get report
	_, err := s.Model.GetReport(ctx, req.Ids[0])
	if err != nil {
		return nil, err
	}

	ids, err := s.Model.UpdateReports(ctx, report, req.Ids)
	if err != nil {
		if status.Code(err) == codes.Unknown {
			return nil, constants.ReportNotFoundError
		}
		return nil, constants.InternalError
	}
	return &pb.UpdateReportsResponse{Data: ids}, nil
}

func (s *UpdateReportsHandler) reqToReport(req *pb.UpdateReportsRequest) *dto.Report {
	report := &dto.Report{
		ID:         req.Data.Id,
		HasSymptom: req.Data.HasSymptom,
	}
	return report
}
