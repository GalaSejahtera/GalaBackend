package covid

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

type GetCovidsHandler struct {
	Model model.IModel
}

func (s *GetCovidsHandler) GetCovids(ctx context.Context, req *pb.GetCovidsRequest) (*pb.GetCovidsResponse, error) {
	var sort *dto.SortData
	var itemsRange *dto.RangeData
	var filter *dto.FilterData

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

	total, covids, err := s.Model.QueryCovids(ctx, sort, itemsRange, filter)
	if err != nil {
		logger.Log.Error("GetCovidsHandler: " + err.Error())
		if status.Code(err) == codes.Unknown {
			return nil, constants.CovidNotFoundError
		}
		return nil, constants.InternalError
	}

	resp, err := s.covidsToResponses(covids)
	if err != nil {
		return nil, err
	}

	resp.Total = total
	return resp, nil
}

func (s *GetCovidsHandler) covidsToResponses(covids []*dto.Covid) (*pb.GetCovidsResponse, error) {
	var resps []*pb.Covid
	for _, covid := range covids {
		resp := &pb.Covid{
			Id:              covid.ID,
			Title:           covid.Title,
			Sid:             covid.SID,
			ImageFeatSingle: covid.ImageFeatSingle,
			Summary:         covid.Summary,
			DatePub2:        covid.DatePub2,
			Content:         covid.Content,
			NewsUrl:         covid.NewsURL,
		}

		resps = append(resps, resp)
	}
	rslt := &pb.GetCovidsResponse{
		Data: resps,
	}

	return rslt, nil
}
