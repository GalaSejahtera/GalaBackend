package kase

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/model"
	"galasejahtera/pkg/utility"
	"github.com/golang/protobuf/ptypes/empty"
)

type GetKasesHandler struct {
	Model model.IModel
}

func (s *GetKasesHandler) GetKases(ctx context.Context, req *empty.Empty) (*pb.GetKasesResponse, error) {
	general := utility.CrawlGeneral()
	resp, err := s.kasesToResponses(general)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GetKasesHandler) kasesToResponses(general *dto.General) (*pb.GetKasesResponse, error) {
	return &pb.GetKasesResponse{
		Data: &pb.General{
			TotalConfirmed: general.TotalConfirmed,
			ActiveCases:    general.ActiveCases,
		},
	}, nil
}
