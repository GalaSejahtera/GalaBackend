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

// CovidDAO ...
type CovidDAO struct {
	client *mongo.Client
}

// InitCovidDAO ...
func InitCovidDAO(client *mongo.Client) ICovidDAO {
	return &CovidDAO{client: client}
}

// Create creates new covid
func (v *CovidDAO) Create(ctx context.Context, covid *dto.Covid) (*dto.Covid, error) {
	// create covid
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Covids)
	if _, err := collection.InsertOne(ctx, covid); err != nil {
		return nil, err
	}
	return covid, nil
}

// Update updates zone
func (v *CovidDAO) Update(ctx context.Context, covid *dto.Covid) (*dto.Covid, error) {
	// update covid
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Covids)
	_, err := collection.UpdateOne(ctx, bson.D{{constants.ID, covid.ID}}, bson.D{
		{"$set", covid},
	})
	if err != nil {
		return nil, err
	}
	return covid, nil
}

// Get gets covid by ID
func (v *CovidDAO) Get(ctx context.Context, id string) (*dto.Covid, error) {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Covids)
	covid := &dto.Covid{}
	if err := collection.FindOne(ctx, bson.D{{constants.ID, id}}).Decode(&covid); err != nil {
		return nil, err
	}
	return covid, nil
}

// BatchGet gets covids by slice of IDs
func (v *CovidDAO) BatchGet(ctx context.Context, ids []string) ([]*dto.Covid, error) {
	var covids []*dto.Covid
	for _, id := range ids {
		covid, err := v.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		covids = append(covids, covid)
	}
	return covids, nil
}

// Query queries covids by sort, range, filter
func (v *CovidDAO) Query(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.Covid, error) {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Covids)

	var cursor *mongo.Cursor
	var err error
	var count int64
	var covids []*dto.Covid

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
			return 1, []*dto.Covid{d}, nil
		}

		// else: do query filter
		if filter.Item == "q" {
			query := bson.M{
				constants.Title: bson.M{
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
		covid := &dto.Covid{}
		if err = cursor.Decode(&covid); err != nil {
			return 0, nil, err
		}
		covids = append(covids, covid)
	}

	return count, covids, nil
}

// Delete deletes covid by ID
func (v *CovidDAO) Delete(ctx context.Context, id string) error {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Covids)
	if _, err := collection.DeleteOne(ctx, bson.D{{constants.ID, id}}); err != nil {
		return err
	}
	return nil
}

// BatchDelete deletes covids by IDs
func (v *CovidDAO) BatchDelete(ctx context.Context, ids []string) ([]string, error) {
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
