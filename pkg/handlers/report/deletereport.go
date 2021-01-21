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

type DeleteReportHandler struct {
	Model model.IModel
}

func (s *DeleteReportHandler) DeleteReport(ctx context.Context, req *pb.DeleteReportRequest) (*pb.DeleteReportResponse, error) {
	rslt, err := s.Model.DeleteReport(ctx, req.Id)
	if err != nil {
		if status.Code(err) == codes.Unknown {
			return nil, constants.ReportNotFoundError
		}
		return nil, constants.InternalError
	}

	resp := s.reportToResp(rslt)
	return resp, nil
}

func (s *DeleteReportHandler) reportToResp(report *dto.Report) *pb.DeleteReportResponse {
	resp := &pb.DeleteReportResponse{
		Data: &pb.Report{
			Id:         report.ID,
			CreatedAt:  report.CreatedAt,
			HasSymptom: report.HasSymptom,
			UserId:     report.UserID,
		},
	}
	return resp
}
