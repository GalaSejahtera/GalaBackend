package dao

import (
	"context"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AuthDAO ...
type AuthDAO struct {
	client *mongo.Client
}

// InitAuthDAO ...
func InitAuthDAO(client *mongo.Client) IAuthDAO {
	return &AuthDAO{client: client}
}

// InitIndex initializes index
func (v *AuthDAO) InitIndex(ctx context.Context) error {
	mod := mongo.IndexModel{
		Keys: bson.M{
			constants.TTL: 1, // index in ascending order
		}, Options: options.Index().SetExpireAfterSeconds(1),
	}

	collection := v.client.Database(constants.SafeWorkout).Collection(constants.AuthTokens)
	_, err := collection.Indexes().CreateOne(ctx, mod)
	return err
}

// Create creates new auth token
func (v *AuthDAO) Create(ctx context.Context, auth *dto.AuthObject) (*dto.AuthObject, error) {
	collection := v.client.Database(constants.SafeWorkout).Collection(constants.AuthTokens)
	if _, err := collection.InsertOne(ctx, auth); err != nil {
		return nil, err
	}
	return auth, nil
}

// Get gets auth token
func (v *AuthDAO) Get(ctx context.Context, token string) (*dto.AuthObject, error) {
	collection := v.client.Database(constants.SafeWorkout).Collection(constants.AuthTokens)
	auth := &dto.AuthObject{}
	if err := collection.FindOne(ctx, bson.D{{constants.Token, token}}).Decode(&auth); err != nil {
		return nil, err
	}
	return auth, nil
}

// Delete deletes user by token
func (v *AuthDAO) Delete(ctx context.Context, token string) error {
	collection := v.client.Database(constants.SafeWorkout).Collection(constants.AuthTokens)
	_, err := collection.DeleteOne(ctx, bson.D{{constants.Token, token}})
	if err != nil {
		return err
	}
	return nil
}

// DeleteByID deletes user by ID
func (v *AuthDAO) DeleteByID(ctx context.Context, id string) error {
	collection := v.client.Database(constants.SafeWorkout).Collection(constants.AuthTokens)

	query := bson.D{{
		"$and",
		bson.A{
			bson.D{{
				constants.UserId,
				bson.D{{
					"$eq",
					id,
				}},
			}},
			bson.D{{
				constants.Type,
				bson.D{{
					"$in",
					[]string{constants.Access, constants.Refresh},
				}},
			}},
		},
	}}

	_, err := collection.DeleteMany(ctx, query)
	if err != nil {
		return err
	}
	return nil
}
