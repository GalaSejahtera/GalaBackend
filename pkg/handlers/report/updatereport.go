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

type UpdateReportHandler struct {
	Model model.IModel
}

func (s *UpdateReportHandler) UpdateReport(ctx context.Context, req *pb.UpdateReportRequest) (*pb.UpdateReportResponse, error) {
	if req.Data == nil {
		return nil, constants.InvalidArgumentError
	}
	report := s.reqToReport(req)

	v, err := s.Model.UpdateReport(ctx, report)
	if err != nil {
		if status.Code(err) == codes.Unknown {
			return nil, constants.ReportNotFoundError
		}
		return nil, constants.InternalError
	}
	resp := s.reportToResp(v)
	return resp, nil
}

func (s *UpdateReportHandler) reqToReport(req *pb.UpdateReportRequest) *dto.Report {
	report := &dto.Report{
		ID:         req.Id,
		HasSymptom: req.Data.HasSymptom,
	}
	return report
}

func (s *UpdateReportHandler) reportToResp(report *dto.Report) *pb.UpdateReportResponse {
	resp := &pb.UpdateReportResponse{
		Data: &pb.Report{
			Id:         report.ID,
			CreatedAt:  report.CreatedAt,
			HasSymptom: report.HasSymptom,
			UserId:     report.UserID,
			Results:    report.Results,
		},
	}
	return resp
}
