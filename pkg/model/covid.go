package model

import (
	"context"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/utility"
)

// GetCovid gets covid by ID
func (m *Model) GetCovid(ctx context.Context, id string) (*dto.Covid, error) {
	utility.CrawlStory(id)
	c, err := m.covidDAO.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Content = utility.CrawlStory(id)
	return c, nil
}

// QueryCovids queries covids by sort, range, filter
func (m *Model) QueryCovids(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.Covid, error) {

	// Commented the block that crawls everything and updates database
	//for i := 0; i < 10; i++ {
	//	covids := utility.CrawlNews(int64(i))
	//	for _, covid := range covids {
	//		_, err := m.covidDAO.Get(ctx, covid.ID)
	//		if err != nil {
	//			// only add into db if not found
	//			_, err = m.covidDAO.Create(ctx, covid)
	//		}
	//	}
	//}

	return m.covidDAO.Query(ctx, sort, itemsRange, filter)
}
