package faq

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "safeworkout/pkg/api"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/dto"
	"safeworkout/pkg/model"
)

type UpdateFaqsHandler struct {
	Model model.IModel
}

func (s *UpdateFaqsHandler) UpdateFaqs(ctx context.Context, req *pb.UpdateFaqsRequest) (*pb.UpdateFaqsResponse, error) {
	if req.Data == nil {
		return nil, constants.InvalidArgumentError
	}
	if len(req.Ids) > 1 {
		return nil, constants.OperationUnsupportedError
	}
	zone := s.reqToFaq(req)

	// get zone
	_, err := s.Model.GetFaq(ctx, req.Ids[0])
	if err != nil {
		return nil, err
	}

	ids, err := s.Model.UpdateFaqs(ctx, zone, req.Ids)
	if err != nil {
		if status.Code(err) == codes.Unknown {
			return nil, constants.FaqNotFoundError
		}
		return nil, constants.InternalError
	}
	return &pb.UpdateFaqsResponse{Data: ids}, nil
}

func (s *UpdateFaqsHandler) reqToFaq(req *pb.UpdateFaqsRequest) *dto.Faq {
	zone := &dto.Faq{
		ID:    req.Data.Id,
		Title: req.Data.Title,
		Desc:  req.Data.Desc,
	}
	return zone
}
