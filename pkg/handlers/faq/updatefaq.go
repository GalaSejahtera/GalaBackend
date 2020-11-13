package faq

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/model"
)

type UpdateFaqHandler struct {
	Model model.IModel
}

func (s *UpdateFaqHandler) UpdateFaq(ctx context.Context, req *pb.UpdateFaqRequest) (*pb.UpdateFaqResponse, error) {
	if req.Data == nil {
		return nil, constants.InvalidArgumentError
	}
	faq := s.reqToFaq(req)

	v, err := s.Model.UpdateFaq(ctx, faq)
	if err != nil {
		if status.Code(err) == codes.Unknown {
			return nil, constants.FaqNotFoundError
		}
		return nil, constants.InternalError
	}
	resp := s.faqToResp(v)
	return resp, nil
}

func (s *UpdateFaqHandler) reqToFaq(req *pb.UpdateFaqRequest) *dto.Faq {
	faq := &dto.Faq{
		ID:    req.Id,
		Title: req.Data.Title,
		Desc:  req.Data.Desc,
	}
	return faq
}

func (s *UpdateFaqHandler) faqToResp(faq *dto.Faq) *pb.UpdateFaqResponse {
	resp := &pb.UpdateFaqResponse{
		Data: &pb.Faq{
			Id:    faq.ID,
			Title: faq.Title,
			Desc:  faq.Desc,
		},
	}
	return resp
}
