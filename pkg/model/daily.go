package model

import (
	"context"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/utility"
	"strconv"
)

// GetDaily gets latest daily
func (m *Model) GetDaily(ctx context.Context) (*dto.Daily, error) {
	_, dailies, err := m.dailyDAO.Query(ctx, &dto.SortData{
		Item:  constants.LastUpdated,
		Order: constants.DESC,
	}, nil, nil)
	if err != nil {
		return nil, err
	}
	if len(dailies) == 0 {
		return &dto.Daily{}, nil
	}
	return dailies[0], nil
}

// UpdateDaily updates dailies (called by scheduler)
func (m *Model) UpdateDailies(ctx context.Context) error {
	daily := utility.CrawlDaily()
	total, _, err := m.dailyDAO.Query(nil, nil, nil, &dto.FilterData{
		Item:  constants.LastUpdated,
		Value: strconv.FormatInt(daily.LastUpdated, 10),
	})
	if err != nil {
		return err
	}
	if total != 0 {
		return nil
	}
	// create new daily
	_, err = m.dailyDAO.Create(ctx, daily)
	return err
}
