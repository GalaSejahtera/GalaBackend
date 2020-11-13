package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
)

// FaqDAO ...
type FaqDAO struct {
	client *mongo.Client
}

// InitFaqDAO ...
func InitFaqDAO(client *mongo.Client) IFaqDAO {
	return &FaqDAO{client: client}
}

// Create creates new faq
func (v *FaqDAO) Create(ctx context.Context, faq *dto.Faq) (*dto.Faq, error) {
	// create faq
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Faqs)
	if _, err := collection.InsertOne(ctx, faq); err != nil {
		return nil, err
	}
	return faq, nil
}

// Update updates zone
func (v *FaqDAO) Update(ctx context.Context, faq *dto.Faq) (*dto.Faq, error) {
	// update faq
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Faqs)
	_, err := collection.UpdateOne(ctx, bson.D{{constants.ID, faq.ID}}, bson.D{
		{"$set", faq},
	})
	if err != nil {
		return nil, err
	}
	return faq, nil
}

// Get gets faq by ID
func (v *FaqDAO) Get(ctx context.Context, id string) (*dto.Faq, error) {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Faqs)
	faq := &dto.Faq{}
	if err := collection.FindOne(ctx, bson.D{{constants.ID, id}}).Decode(&faq); err != nil {
		return nil, err
	}
	return faq, nil
}

// BatchGet gets faqs by slice of IDs
func (v *FaqDAO) BatchGet(ctx context.Context, ids []string) ([]*dto.Faq, error) {
	var faqs []*dto.Faq
	for _, id := range ids {
		faq, err := v.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		faqs = append(faqs, faq)
	}
	return faqs, nil
}

// Query queries faqs by sort, range, filter
func (v *FaqDAO) Query(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.Faq, error) {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Faqs)

	var cursor *mongo.Cursor
	var err error
	var count int64
	var faqs []*dto.Faq

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
			return 1, []*dto.Faq{d}, nil
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
		faq := &dto.Faq{}
		if err = cursor.Decode(&faq); err != nil {
			return 0, nil, err
		}
		faqs = append(faqs, faq)
	}

	return count, faqs, nil
}

// Delete deletes faq by ID
func (v *FaqDAO) Delete(ctx context.Context, id string) error {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Faqs)
	if _, err := collection.DeleteOne(ctx, bson.D{{constants.ID, id}}); err != nil {
		return err
	}
	return nil
}

// BatchDelete deletes faqs by IDs
func (v *FaqDAO) BatchDelete(ctx context.Context, ids []string) ([]string, error) {
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
