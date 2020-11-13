package faq

import (
	"context"
	pb "safeworkout/pkg/api"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/dto"
	"safeworkout/pkg/logger"
	"safeworkout/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetFaqHandler struct {
	Model model.IModel
}

func (s *GetFaqHandler) GetFaq(ctx context.Context, req *pb.GetFaqRequest) (*pb.GetFaqResponse, error) {
	faq, err := s.Model.GetFaq(ctx, req.Id)
	if err != nil {
		logger.Log.Error("GetFaqHandler: " + err.Error())
		if status.Code(err) == codes.Unknown {
			return nil, constants.FaqNotFoundError
		}
		return nil, constants.InternalError
	}

	resp, err := s.faqToResponse(faq)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GetFaqHandler) faqToResponse(faq *dto.Faq) (*pb.GetFaqResponse, error) {
	resp := &pb.GetFaqResponse{
		Data: &pb.Faq{
			Id:    faq.ID,
			Title: faq.Title,
			Desc:  faq.Desc,
		},
	}

	return resp, nil
}
