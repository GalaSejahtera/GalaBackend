package zone

import (
	"context"
	pb "galasejahtera/pkg/api"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/model"
	"sort"
)

type GetRecentZonesByUserHandler struct {
	Model model.IModel
}

func (s *GetRecentZonesByUserHandler) GetRecentZonesByUser(ctx context.Context, req *pb.GetZonesRequest) (*pb.GetZonesResponse, error) {

	zones, err := s.Model.GetRecentZonesByUser(ctx, req.FilterValue)
	if err != nil {
		return nil, err
	}

	if req.Item != "" && req.Order != "" {
		s.sortZones(zones, req.Item, req.Order)
	}

	// get total
	total := int64(len(zones))

	index := int64(0)
	var rslt []*dto.Zone

	for _, zone := range zones {
		if index < req.From {
			index += 1
			continue
		}
		if req.To != 0 && index > req.To {
			index += 1
			continue
		}
		index += 1
		rslt = append(rslt, zone)
	}

	resp := s.zonesToResponses(rslt)
	resp.Total = total

	return resp, nil
}

func (s *GetRecentZonesByUserHandler) zonesToResponses(zones []*dto.Zone) *pb.GetZonesResponse {
	var resps []*pb.Zone
	for _, zone := range zones {
		resp := &pb.Zone{
			Id:       zone.ID,
			Name:     zone.Name,
			Lat:      zone.Lat,
			Long:     zone.Long,
			Type:     zone.Type,
			Radius:   zone.Radius,
			Capacity: zone.Capacity,
			Time:     zone.Time,
			Risk:     zone.Risk,
		}
		resps = append(resps, resp)
	}
	rslt := &pb.GetZonesResponse{
		Data: resps,
	}

	return rslt
}

func (s *GetRecentZonesByUserHandler) sortZones(zones []*dto.Zone, field string, order string) {
	switch field {
	case "id":
		sort.Slice(zones, func(i, j int) bool {
			return zones[i].ID < zones[j].ID
		})
	case "name":
		sort.Slice(zones, func(i, j int) bool {
			return zones[i].Name < zones[j].Name
		})
	case "type":
		sort.Slice(zones, func(i, j int) bool {
			return zones[i].Type < zones[j].Type
		})
	case "radius":
		sort.Slice(zones, func(i, j int) bool {
			return zones[i].Radius < zones[j].Radius
		})
	case "capacity":
		sort.Slice(zones, func(i, j int) bool {
			return zones[i].Capacity < zones[j].Capacity
		})
	case "time":
		sort.Slice(zones, func(i, j int) bool {
			return zones[i].Time < zones[j].Time
		})
	case "risk":
		sort.Slice(zones, func(i, j int) bool {
			return zones[i].Risk < zones[j].Risk
		})
	default:
	}

	if order == "DESC" {
		// reverse slice
		for i, j := 0, len(zones)-1; i < j; i, j = i+1, j-1 {
			zones[i], zones[j] = zones[j], zones[i]
		}
	}
}
