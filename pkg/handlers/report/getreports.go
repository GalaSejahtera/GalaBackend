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

type GetReportsHandler struct {
	Model model.IModel
}

func (s *GetReportsHandler) GetReports(ctx context.Context, req *pb.GetReportsRequest) (*pb.GetReportsResponse, error) {
	var sort *dto.SortData
	var itemsRange *dto.RangeData
	var filter *dto.FilterData

	// If the request is batch get, call batch get model
	if len(req.Ids) > 0 {
		reports, err := s.Model.BatchGetReports(ctx, req.Ids)
		if err != nil {
			if status.Code(err) == codes.Unknown {
				return nil, constants.ReportNotFoundError
			}
			return nil, constants.InternalError
		}
		resp, err := s.reportsToResponses(reports)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}

	if req.Item != "" && req.Order != "" {
		sort = &dto.SortData{
			Item:  req.Item,
			Order: req.Order,
		}
	}

	if req.To != 0 {
		itemsRange = &dto.RangeData{
			From: int(req.From),
			To:   int(req.To),
		}
	}

	if req.FilterItem != "" && req.FilterValue != "" {
		filter = &dto.FilterData{
			Item:  req.FilterItem,
			Value: req.FilterValue,
		}
	}

	total, reports, err := s.Model.QueryReports(ctx, sort, itemsRange, filter)
	if err != nil {
		logger.Log.Error("GetReportsHandler: " + err.Error())
		if status.Code(err) == codes.Unknown {
			return nil, constants.ReportNotFoundError
		}
		return nil, constants.InternalError
	}

	resp, err := s.reportsToResponses(reports)
	if err != nil {
		return nil, err
	}

	resp.Total = total
	return resp, nil
}

func (s *GetReportsHandler) reportsToResponses(reports []*dto.Report) (*pb.GetReportsResponse, error) {
	var resps []*pb.Report
	for _, report := range reports {
		resp := &pb.Report{
			Id:         report.ID,
			CreatedAt:  report.CreatedAt,
			HasSymptom: report.HasSymptom,
			UserId:     report.UserID,
		}

		resps = append(resps, resp)
	}
	rslt := &pb.GetReportsResponse{
		Data: resps,
	}

	return rslt, nil
}
