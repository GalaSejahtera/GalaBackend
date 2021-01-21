package model

import (
	"context"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateReport creates new report
func (m *Model) CreateReport(ctx context.Context, report *dto.Report) (*dto.Report, error) {

	// check if Report exist
	_, err := m.reportDAO.Get(ctx, report.ID)

	// only can create report if not found
	if err != nil && status.Code(err) == codes.Unknown {
		return m.reportDAO.Create(ctx, report)
	}

	if err != nil {
		return nil, err
	}

	return nil, status.Error(codes.AlreadyExists, "Report already exist!")
}

// UpdateReport updates report
func (m *Model) UpdateReport(ctx context.Context, report *dto.Report) (*dto.Report, error) {

	// check if report exists
	f, err := m.reportDAO.Get(ctx, report.ID)
	if err != nil {
		return nil, err
	}

	// patch report
	f.HasSymptom = report.HasSymptom

	_, err = m.reportDAO.Update(ctx, f)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// UpdateReports update reports
func (m *Model) UpdateReports(ctx context.Context, report *dto.Report, ids []string) ([]string, error) {
	if len(ids) > 1 {
		return nil, constants.OperationUnsupportedError
	}
	report.ID = ids[0]
	u, err := m.UpdateReport(ctx, report)
	if err != nil {
		return nil, err
	}
	return []string{u.ID}, err
}

// GetReport gets report by ID
func (m *Model) GetReport(ctx context.Context, id string) (*dto.Report, error) {
	return m.reportDAO.Get(ctx, id)
}

// BatchGetReports get reports by slice of IDs
func (m *Model) BatchGetReports(ctx context.Context, ids []string) ([]*dto.Report, error) {
	return m.reportDAO.BatchGet(ctx, ids)
}

// QueryReports queries reports by sort, range, filter
func (m *Model) QueryReports(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.Report, error) {
	return m.reportDAO.Query(ctx, sort, itemsRange, filter)
}

// DeleteReport deletes report by ID
func (m *Model) DeleteReport(ctx context.Context, id string) (*dto.Report, error) {
	// check if report exist
	u, err := m.reportDAO.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = m.reportDAO.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// DeleteReports delete reports by IDs
func (m *Model) DeleteReports(ctx context.Context, ids []string) ([]string, error) {
	var deletedIDs []string
	for _, id := range ids {
		u, err := m.DeleteReport(ctx, id)
		if err != nil {
			return nil, err
		}
		deletedIDs = append(deletedIDs, u.ID)
	}

	return deletedIDs, nil
}
