package kase

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/model"
	"galasejahtera/pkg/utility"
	"github.com/golang/protobuf/ptypes/empty"
	"time"
)

type GetRecentKasesHandler struct {
	Model model.IModel
}

func (s *GetRecentKasesHandler) GetRecentKases(ctx context.Context, req *empty.Empty) (*pb.GetRecentKasesResponse, error) {
	// get date range
	now := utility.MalaysiaTime(time.Now())
	t2 := utility.TimeToDateStringWithDash(now)
	t1 := utility.TimeToDateStringWithDash(now.Add(-24 * 8 * time.Hour))

	kases := utility.CrawlCasesByDate(t1, t2)
	resp, err := s.kasesToResponses(kases)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GetRecentKasesHandler) kasesToResponses(kases []*dto.Kase) (*pb.GetRecentKasesResponse, error) {
	data := []*pb.Kase{}
	for _, k := range kases {
		data = append(data, &pb.Kase{
			LastUpdated:   k.LastUpdated,
			NewDeaths:     k.NewDeaths,
			NewInfections: k.NewInfections,
			NewRecovered:  k.NewRecovered,
		})
	}

	return &pb.GetRecentKasesResponse{
		Data: data,
	}, nil
}
