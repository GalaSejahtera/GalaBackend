package faq

import (
	"context"
	"github.com/twinj/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/model"
)

type CreateFaqHandler struct {
	Model model.IModel
}

func (s *CreateFaqHandler) CreateFaq(ctx context.Context, req *pb.CreateFaqRequest) (*pb.CreateFaqResponse, error) {
	if req.Data == nil {
		return nil, constants.InvalidArgumentError
	}
	faq := &dto.Faq{
		ID:    uuid.NewV4().String(),
		Title: req.Data.Title,
		Desc:  req.Data.Desc,
	}

	rslt, err := s.Model.CreateFaq(ctx, faq)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return nil, constants.FaqAlreadyExistError
		}
		return nil, constants.InternalError
	}
	resp := s.faqToResp(rslt)
	return resp, nil
}

func (s *CreateFaqHandler) faqToResp(faq *dto.Faq) *pb.CreateFaqResponse {
	return &pb.CreateFaqResponse{
		Data: &pb.Faq{
			Id:    faq.ID,
			Title: faq.Title,
			Desc:  faq.Desc,
		},
	}
}
