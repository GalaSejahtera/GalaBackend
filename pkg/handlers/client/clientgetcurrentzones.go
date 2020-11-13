package client

import (
	"context"
	pb "safeworkout/pkg/api"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/dto"
	"safeworkout/pkg/logger"
	"safeworkout/pkg/model"
	"safeworkout/pkg/utility"
	"time"

	"github.com/twinj/uuid"
)

type GetCurrentZonesHandler struct {
	Model model.IModel
}

func (s *GetCurrentZonesHandler) GetCurrentZones(ctx context.Context, req *pb.ClientGetCurrentZonesRequest) (*pb.ClientGetCurrentZonesResponse, error) {
	if req.User == nil {
		return nil, constants.InvalidArgumentError
	}

	u := &dto.User{
		ID:          req.User.Id,
		IsActive:    true,
		LastUpdated: utility.TimeToMilli(utility.MalaysiaTime(time.Now())),
		Lat:         req.User.Lat,
		Long:        req.User.Long,
	}

	// update user location, isActive and lastUpdated
	_, err := s.Model.ClientUpdateUser(ctx, u)
	if err != nil {
		logger.Log.Error("ClientGetCurrentZonesHandler: " + err.Error())
		return nil, err
	}

	// get user current zone and sub zones
	zone, subZones, err := s.Model.GetZonesByUser(ctx, u.ID)
	if err != nil {
		logger.Log.Error("ClientGetCurrentZonesHandler: " + err.Error())
		return nil, err
	}

	// create activity
	_, err = s.Model.CreateActivity(ctx, &dto.Activity{
		ID:     uuid.NewV4().String(),
		ZoneID: zone.ID,
		UserID: u.ID,
		Time:   u.LastUpdated,
	})
	if err != nil {
		logger.Log.Warn("ClientGetCurrentZonesHandler: " + err.Error())
	}

	resp, err := s.zonesToResponses(zone, subZones)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GetCurrentZonesHandler) zonesToResponses(zone *dto.Zone, subZones []*dto.Zone) (*pb.ClientGetCurrentZonesResponse, error) {
	zonePb := &pb.Zone{
		Id:                 zone.ID,
		Name:               zone.Name,
		Lat:                zone.Lat,
		Long:               zone.Long,
		Type:               zone.Type,
		Radius:             zone.Radius,
		Capacity:           zone.Capacity,
		UsersWithin:        zone.UsersWithin,
		IsCapacityExceeded: zone.IsCapacityExceeded,
		Risk:               zone.Risk,
	}

	subZonesPb := []*pb.Zone{}
	for _, z := range subZones {
		zPb := &pb.Zone{
			Id:                 z.ID,
			Name:               z.Name,
			Lat:                z.Lat,
			Long:               z.Long,
			Type:               z.Type,
			Radius:             z.Radius,
			Capacity:           z.Capacity,
			UsersWithin:        z.UsersWithin,
			IsCapacityExceeded: z.IsCapacityExceeded,
			Risk:               z.Risk,
		}
		subZonesPb = append(subZonesPb, zPb)
	}

	rslt := &pb.ClientGetCurrentZonesResponse{
		Zone:     zonePb,
		SubZones: subZonesPb,
	}

	return rslt, nil
}
