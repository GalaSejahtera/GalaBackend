package faq

import (
	"context"
	pb "safeworkout/pkg/api"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/dto"
	"safeworkout/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteFaqHandler struct {
	Model model.IModel
}

func (s *DeleteFaqHandler) DeleteFaq(ctx context.Context, req *pb.DeleteFaqRequest) (*pb.DeleteFaqResponse, error) {
	rslt, err := s.Model.DeleteFaq(ctx, req.Id)
	if err != nil {
		if status.Code(err) == codes.Unknown {
			return nil, constants.FaqNotFoundError
		}
		return nil, constants.InternalError
	}

	resp := s.faqToResp(rslt)
	return resp, nil
}

func (s *DeleteFaqHandler) faqToResp(faq *dto.Faq) *pb.DeleteFaqResponse {
	resp := &pb.DeleteFaqResponse{
		Data: &pb.Faq{
			Id:    faq.ID,
			Title: faq.Title,
			Desc:  faq.Desc,
		},
	}
	return resp
}
