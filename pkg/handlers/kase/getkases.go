package kase

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/model"
	"github.com/golang/protobuf/ptypes/empty"
)

type GetKasesHandler struct {
	Model model.IModel
}

func (s *GetKasesHandler) GetKases(ctx context.Context, req *empty.Empty) (*pb.GetKasesResponse, error) {
	resp, err := s.kasesToResponses()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GetKasesHandler) kasesToResponses() (*pb.GetKasesResponse, error) {
	rslt := map[string]int64{
		"ok": 1,
	}

	return &pb.GetKasesResponse{
		Data: rslt,
	}, nil
}
