package model

import (
	"context"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/dto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateFaq creates new faq
func (m *Model) CreateFaq(ctx context.Context, faq *dto.Faq) (*dto.Faq, error) {

	// check if Faq exist
	_, err := m.faqDAO.Get(ctx, faq.ID)

	// only can create faq if not found
	if err != nil && status.Code(err) == codes.Unknown {
		return m.faqDAO.Create(ctx, faq)
	}

	if err != nil {
		return nil, err
	}

	return nil, status.Error(codes.AlreadyExists, "Faq already exist!")
}

// UpdateFaq updates faq
func (m *Model) UpdateFaq(ctx context.Context, faq *dto.Faq) (*dto.Faq, error) {

	// check if faq exists
	f, err := m.faqDAO.Get(ctx, faq.ID)
	if err != nil {
		return nil, err
	}

	// patch faq
	f.Title = faq.Title
	f.Desc = faq.Desc

	_, err = m.faqDAO.Update(ctx, f)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// UpdateFaqs update faqs
func (m *Model) UpdateFaqs(ctx context.Context, faq *dto.Faq, ids []string) ([]string, error) {
	if len(ids) > 1 {
		return nil, constants.OperationUnsupportedError
	}
	faq.ID = ids[0]
	u, err := m.UpdateFaq(ctx, faq)
	if err != nil {
		return nil, err
	}
	return []string{u.ID}, err
}

// GetFaq gets faq by ID
func (m *Model) GetFaq(ctx context.Context, id string) (*dto.Faq, error) {
	return m.faqDAO.Get(ctx, id)
}

// BatchGetFaqs get faqs by slice of IDs
func (m *Model) BatchGetFaqs(ctx context.Context, ids []string) ([]*dto.Faq, error) {
	return m.faqDAO.BatchGet(ctx, ids)
}

// QueryFaqs queries faqs by sort, range, filter
func (m *Model) QueryFaqs(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.Faq, error) {
	return m.faqDAO.Query(ctx, sort, itemsRange, filter)
}

// DeleteFaq deletes faq by ID
func (m *Model) DeleteFaq(ctx context.Context, id string) (*dto.Faq, error) {
	// check if faq exist
	u, err := m.faqDAO.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = m.faqDAO.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// DeleteFaqs delete faqs by IDs
func (m *Model) DeleteFaqs(ctx context.Context, ids []string) ([]string, error) {
	var deletedIDs []string
	for _, id := range ids {
		u, err := m.DeleteFaq(ctx, id)
		if err != nil {
			return nil, err
		}
		deletedIDs = append(deletedIDs, u.ID)
	}

	return deletedIDs, nil
}
