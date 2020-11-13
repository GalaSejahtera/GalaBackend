package faq

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/logger"
	"galasejahtera/pkg/model"
)

type GetFaqsHandler struct {
	Model model.IModel
}

func (s *GetFaqsHandler) GetFaqs(ctx context.Context, req *pb.GetFaqsRequest) (*pb.GetFaqsResponse, error) {
	var sort *dto.SortData
	var itemsRange *dto.RangeData
	var filter *dto.FilterData

	// If the request is batch get, call batch get model
	if len(req.Ids) > 0 {
		faqs, err := s.Model.BatchGetFaqs(ctx, req.Ids)
		if err != nil {
			if status.Code(err) == codes.Unknown {
				return nil, constants.FaqNotFoundError
			}
			return nil, constants.InternalError
		}
		resp, err := s.faqsToResponses(faqs)
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

	total, faqs, err := s.Model.QueryFaqs(ctx, sort, itemsRange, filter)
	if err != nil {
		logger.Log.Error("GetFaqsHandler: " + err.Error())
		if status.Code(err) == codes.Unknown {
			return nil, constants.FaqNotFoundError
		}
		return nil, constants.InternalError
	}

	resp, err := s.faqsToResponses(faqs)
	if err != nil {
		return nil, err
	}

	resp.Total = total
	return resp, nil
}

func (s *GetFaqsHandler) faqsToResponses(faqs []*dto.Faq) (*pb.GetFaqsResponse, error) {
	var resps []*pb.Faq
	for _, faq := range faqs {
		resp := &pb.Faq{
			Id:    faq.ID,
			Title: faq.Title,
			Desc:  faq.Desc,
		}

		resps = append(resps, resp)
	}
	rslt := &pb.GetFaqsResponse{
		Data: resps,
	}

	return rslt, nil
}
