package faq

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteFaqsHandler struct {
	Model model.IModel
}

func (s *DeleteFaqsHandler) DeleteFaqs(ctx context.Context, req *pb.DeleteFaqsRequest) (*pb.DeleteFaqsResponse, error) {
	var ids []string

	// remove faqs
	for _, id := range req.Ids {
		u, err := s.Model.DeleteFaq(ctx, id)
		if err != nil {
			if status.Code(err) == codes.Unknown {
				return nil, constants.FaqNotFoundError
			}
			return nil, constants.InternalError
		}

		// add faq into deleted faq IDs
		ids = append(ids, u.ID)
	}

	return &pb.DeleteFaqsResponse{Data: ids}, nil
}
