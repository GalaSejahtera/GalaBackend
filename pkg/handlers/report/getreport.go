package report

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

type GetReportHandler struct {
	Model model.IModel
}

func (s *GetReportHandler) GetReport(ctx context.Context, req *pb.GetReportRequest) (*pb.GetReportResponse, error) {
	report, err := s.Model.GetReport(ctx, req.Id)
	if err != nil {
		logger.Log.Error("GetReportHandler: " + err.Error())
		if status.Code(err) == codes.Unknown {
			return nil, constants.ReportNotFoundError
		}
		return nil, constants.InternalError
	}

	resp, err := s.reportToResponse(report)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GetReportHandler) reportToResponse(report *dto.Report) (*pb.GetReportResponse, error) {
	resp := &pb.GetReportResponse{
		Data: &pb.Report{
			Id:         report.ID,
			CreatedAt:  report.CreatedAt,
			HasSymptom: report.HasSymptom,
			UserId:     report.UserID,
		},
	}

	return resp, nil
}
