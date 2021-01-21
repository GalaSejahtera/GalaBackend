package report

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/model"
	"galasejahtera/pkg/utility"
	"github.com/twinj/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type CreateReportHandler struct {
	Model model.IModel
}

func (s *CreateReportHandler) CreateReport(ctx context.Context, req *pb.CreateReportRequest) (*pb.CreateReportResponse, error) {
	if req.Data == nil {
		return nil, constants.InvalidArgumentError
	}
	report := &dto.Report{
		ID:         uuid.NewV4().String(),
		CreatedAt:  utility.TimeToMilli(utility.MalaysiaTime(time.Now())),
		HasSymptom: req.Data.HasSymptom,
		UserID:     req.Data.UserId,
	}

	rslt, err := s.Model.CreateReport(ctx, report)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return nil, constants.ReportAlreadyExistError
		}
		return nil, constants.InternalError
	}
	resp := s.reportToResp(rslt)
	return resp, nil
}

func (s *CreateReportHandler) reportToResp(report *dto.Report) *pb.CreateReportResponse {
	return &pb.CreateReportResponse{
		Data: &pb.Report{
			Id:         report.ID,
			CreatedAt:  report.CreatedAt,
			HasSymptom: report.HasSymptom,
			UserId:     report.UserID,
		},
	}
}
