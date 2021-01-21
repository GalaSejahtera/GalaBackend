package daily

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/logger"
	"galasejahtera/pkg/model"
	"galasejahtera/pkg/utility"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetDailyHandler struct {
	Model model.IModel
}

func (s *GetDailyHandler) GetDistrict(ctx context.Context, req *pb.GetDistrictRequest) (*pb.GetDistrictResponse, error) {
	daily, err := s.Model.GetDaily(ctx)
	if err != nil {
		logger.Log.Error("GetDailyHandler: " + err.Error())
		if status.Code(err) == codes.Unknown {
			return nil, constants.DailyNotFoundError
		}
		return nil, constants.InternalError
	}

	d := &dto.District{}
	for _, state := range daily.States {
		for _, district := range state.Districts {
			if utility.NormalizePlace(district.Name) == utility.NormalizePlace(req.Id) {
				d = district
				break
			}
		}
		if utility.NormalizePlace(state.Name) == utility.NormalizePlace(req.Id) {
			d.Name = state.Name
			d.Total = state.Total
			break
		}
	}

	resp, err := s.districtToResponse(d)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GetDailyHandler) districtToResponse(district *dto.District) (*pb.GetDistrictResponse, error) {
	resp := &pb.GetDistrictResponse{
		Data: &pb.District{
			Name:  district.Name,
			Total: district.Total,
		},
	}

	return resp, nil
}
