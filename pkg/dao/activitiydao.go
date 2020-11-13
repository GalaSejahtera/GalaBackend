package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/utility"
	"time"
)

// ActivityDAO ...
type ActivityDAO struct {
	client *mongo.Client
}

// InitActivityDAO ...
func InitActivityDAO(client *mongo.Client) IActivityDAO {
	return &ActivityDAO{client: client}
}

// Create creates new activity
func (v *ActivityDAO) Create(ctx context.Context, activity *dto.Activity) (*dto.Activity, error) {
	activity.TTL = utility.MilliToTime(time.Now().Add(time.Hour*24).Unix()*1000 - 1000)
	// create activity
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Activities)
	if _, err := collection.InsertOne(ctx, activity); err != nil {
		return nil, err
	}
	return activity, nil
}

// Get gets activity by ID
func (v *ActivityDAO) Get(ctx context.Context, id string) (*dto.Activity, error) {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Activities)
	activity := &dto.Activity{}
	if err := collection.FindOne(ctx, bson.D{{constants.ID, id}}).Decode(&activity); err != nil {
		return nil, err
	}
	return activity, nil
}

// Query queries activities by sort, range, filter
func (v *ActivityDAO) Query(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.Activity, error) {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Activities)

	var cursor *mongo.Cursor
	var err error
	var count int64
	var activities []*dto.Activity

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
			return 1, []*dto.Activity{d}, nil
		}

		// else: do query filter
		if filter.Item == "q" {
			query := bson.M{
				constants.ID: bson.M{
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
		activity := &dto.Activity{}
		if err = cursor.Decode(&activity); err != nil {
			return 0, nil, err
		}
		activities = append(activities, activity)
	}

	return count, activities, nil
}

// Delete deletes activity by ID
func (v *ActivityDAO) Delete(ctx context.Context, id string) error {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Activities)
	if _, err := collection.DeleteOne(ctx, bson.D{{constants.ID, id}}); err != nil {
		return err
	}
	return nil
}

// BatchDelete deletes activities by IDs
func (v *ActivityDAO) BatchDelete(ctx context.Context, ids []string) ([]string, error) {
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
