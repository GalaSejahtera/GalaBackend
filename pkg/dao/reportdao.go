package dao

import (
	"context"
	"fmt"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ReportDAO ...
type ReportDAO struct {
	client *mongo.Client
}

// InitReportDAO ...
func InitReportDAO(client *mongo.Client) IReportDAO {
	return &ReportDAO{client: client}
}

// Create creates new report
func (v *ReportDAO) Create(ctx context.Context, report *dto.Report) (*dto.Report, error) {
	// create report
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Reports)
	if _, err := collection.InsertOne(ctx, report); err != nil {
		return nil, err
	}
	return report, nil
}

// Update updates zone
func (v *ReportDAO) Update(ctx context.Context, report *dto.Report) (*dto.Report, error) {
	// update report
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Reports)
	_, err := collection.UpdateOne(ctx, bson.D{{constants.ID, report.ID}}, bson.D{
		{"$set", report},
	})
	if err != nil {
		return nil, err
	}
	return report, nil
}

// Get gets report by ID
func (v *ReportDAO) Get(ctx context.Context, id string) (*dto.Report, error) {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Reports)
	report := &dto.Report{}
	if err := collection.FindOne(ctx, bson.D{{constants.ID, id}}).Decode(&report); err != nil {
		return nil, err
	}
	return report, nil
}

// BatchGet gets reports by slice of IDs
func (v *ReportDAO) BatchGet(ctx context.Context, ids []string) ([]*dto.Report, error) {
	var reports []*dto.Report
	for _, id := range ids {
		report, err := v.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, nil
}

// Query queries reports by sort, range, filter
func (v *ReportDAO) Query(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.Report, error) {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Reports)

	var cursor *mongo.Cursor
	var err error
	var count int64
	var reports []*dto.Report

	findOptions := options.Find()
	// set range
	if itemsRange != nil {
		findOptions.SetSkip(int64(itemsRange.From))
		findOptions.SetLimit(int64(itemsRange.To + 1 - itemsRange.From))
	}

	// set sorter
	if sort != nil {
		order := 1
		if sort.Order == constants.DESC {
			order = -1
		}
		findOptions.SetSort(bson.D{{sort.Item, order}})
	}

	// set filter
	if filter != nil {
		// special case: if filter item is primary key we can directly call Get
		if filter.Item == "id" {
			d, err := v.Get(ctx, filter.Value)
			if err != nil {
				return 0, nil, err
			}
			return 1, []*dto.Report{d}, nil
		}

		// else: do query filter
		if filter.Item == "q" {
			query := bson.M{
				constants.UserId: bson.M{
					"$regex":   fmt.Sprintf("%s.*", filter.Value),
					"$options": "i",
				},
			}
			cursor, err = collection.Find(ctx, query, findOptions)
			if err != nil {
				return 0, nil, err
			}
			count, err = collection.CountDocuments(ctx, query)
			if err != nil {
				return 0, nil, err
			}
		} else {
			cursor, err = collection.Find(
				ctx, bson.D{
					{filter.Item, filter.Value},
				}, findOptions,
			)
			if err != nil {
				return 0, nil, err
			}
			count, err = collection.CountDocuments(ctx, bson.D{
				{filter.Item, filter.Value},
			})
			if err != nil {
				return 0, nil, err
			}
		}
	} else {
		cursor, err = collection.Find(ctx, bson.D{{}}, findOptions)
		if err != nil {
			return 0, nil, err
		}
		count, err = collection.CountDocuments(ctx, bson.D{{}})
		if err != nil {
			return 0, nil, err
		}
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		report := &dto.Report{}
		if err = cursor.Decode(&report); err != nil {
			return 0, nil, err
		}
		reports = append(reports, report)
	}

	return count, reports, nil
}

// Delete deletes report by ID
func (v *ReportDAO) Delete(ctx context.Context, id string) error {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Reports)
	if _, err := collection.DeleteOne(ctx, bson.D{{constants.ID, id}}); err != nil {
		return err
	}
	return nil
}

// BatchDelete deletes reports by IDs
func (v *ReportDAO) BatchDelete(ctx context.Context, ids []string) ([]string, error) {
	var deletedIDs []string
	for _, id := range ids {
		err := v.Delete(ctx, id)
		if err != nil {
			return nil, err
		}
		deletedIDs = append(deletedIDs, id)
	}
	return deletedIDs, nil
}
