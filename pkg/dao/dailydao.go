package dao

import (
	"context"
	"fmt"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

// DailyDAO ...
type DailyDAO struct {
	client *mongo.Client
}

// InitDailyDAO ...
func InitDailyDAO(client *mongo.Client) IDailyDAO {
	return &DailyDAO{client: client}
}

// Create creates new daily
func (v *DailyDAO) Create(ctx context.Context, daily *dto.Daily) (*dto.Daily, error) {
	// create daily
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Dailies)
	if _, err := collection.InsertOne(ctx, daily); err != nil {
		return nil, err
	}
	return daily, nil
}

// Get gets daily by ID
func (v *DailyDAO) Get(ctx context.Context, id string) (*dto.Daily, error) {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Dailies)
	daily := &dto.Daily{}
	if err := collection.FindOne(ctx, bson.D{{constants.ID, id}}).Decode(&daily); err != nil {
		return nil, err
	}
	return daily, nil
}

// BatchGet gets dailies by slice of IDs
func (v *DailyDAO) BatchGet(ctx context.Context, ids []string) ([]*dto.Daily, error) {
	var dailies []*dto.Daily
	for _, id := range ids {
		daily, err := v.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		dailies = append(dailies, daily)
	}
	return dailies, nil
}

// Query queries dailies by sort, range, filter
func (v *DailyDAO) Query(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.Daily, error) {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Dailies)

	var cursor *mongo.Cursor
	var err error
	var count int64
	var dailies []*dto.Daily

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
			return 1, []*dto.Daily{d}, nil
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
		} else if filter.Item == constants.LastUpdated {
			n, err := strconv.ParseInt(filter.Value, 10, 64)
			if err != nil {
				return 0, nil, err
			}
			cursor, err = collection.Find(
				ctx, bson.D{
					{filter.Item, n},
				}, findOptions,
			)
			if err != nil {
				return 0, nil, err
			}
			count, err = collection.CountDocuments(ctx, bson.D{
				{filter.Item, n},
			})
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
		daily := &dto.Daily{}
		if err = cursor.Decode(&daily); err != nil {
			return 0, nil, err
		}
		dailies = append(dailies, daily)
	}

	return count, dailies, nil
}

// Delete deletes daily by ID
func (v *DailyDAO) Delete(ctx context.Context, id string) error {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Dailies)
	if _, err := collection.DeleteOne(ctx, bson.D{{constants.ID, id}}); err != nil {
		return err
	}
	return nil
}

// BatchDelete deletes dailies by IDs
func (v *DailyDAO) BatchDelete(ctx context.Context, ids []string) ([]string, error) {
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

// Query queries dailies by time range
func (v *DailyDAO) QueryByTimeRange(ctx context.Context, startTime int64, endTime int64) (int64, []*dto.Daily, error) {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Dailies)

	var cursor *mongo.Cursor
	var err error
	var count int64
	var dailies []*dto.Daily

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{constants.LastUpdated, -1}})
	query := bson.M{
		constants.LastUpdated: bson.M{
			"$gte": startTime,
			"$lte": endTime,
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
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		daily := &dto.Daily{}
		if err = cursor.Decode(&daily); err != nil {
			return 0, nil, err
		}
		dailies = append(dailies, daily)
	}

	return count, dailies, nil
}
