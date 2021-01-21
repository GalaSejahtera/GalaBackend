package covid

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/logger"
	"galasejahtera/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetCovidHandler struct {
	Model model.IModel
}

func (s *GetCovidHandler) GetCovid(ctx context.Context, req *pb.GetCovidRequest) (*pb.GetCovidResponse, error) {
	covid, err := s.Model.GetCovid(ctx, req.Id)
	if err != nil {
		logger.Log.Error("GetCovidHandler: " + err.Error())
		if status.Code(err) == codes.Unknown {
			return nil, constants.CovidNotFoundError
		}
		return nil, constants.InternalError
	}

	resp, err := s.covidToResponse(covid)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GetCovidHandler) covidToResponse(covid *dto.Covid) (*pb.GetCovidResponse, error) {
	resp := &pb.GetCovidResponse{
		Data: &pb.Covid{
			Id:              covid.ID,
			Title:           covid.Title,
			Sid:             covid.SID,
			ImageFeatSingle: covid.ImageFeatSingle,
			Summary:         covid.Summary,
			DatePub2:        covid.DatePub2,
			Content:         covid.Content,
			NewsUrl:         covid.NewsURL,
		},
	}

	return resp, nil
}
