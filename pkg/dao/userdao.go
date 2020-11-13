package dao

import (
	"context"
	"fmt"
	"galasejahtera/pkg/constants"
	"galasejahtera/pkg/dto"
	"galasejahtera/pkg/utility"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserDAO ...
type UserDAO struct {
	client *mongo.Client
}

// InitUserDAO ...
func InitUserDAO(client *mongo.Client) IUserDAO {
	return &UserDAO{client: client}
}

// Create creates new user
func (v *UserDAO) Create(ctx context.Context, user *dto.User) (*dto.User, error) {

	// hash password
	var err error
	user.Password, err = utility.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	// create location
	location := &dto.Location{
		Type:        "Point",
		Coordinates: []float64{user.Long, user.Lat},
	}
	user.Location = location

	// create user
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Users)
	if _, err := collection.InsertOne(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// Get gets user by ID
func (v *UserDAO) Get(ctx context.Context, id string) (*dto.User, error) {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Users)
	user := &dto.User{}
	if err := collection.FindOne(ctx, bson.D{{constants.ID, id}}).Decode(&user); err != nil {
		return nil, err
	}

	if user.Location != nil && len(user.Location.Coordinates) == 2 {
		user.Long = user.Location.Coordinates[0]
		user.Lat = user.Location.Coordinates[1]
	}

	return user, nil
}

// BatchGet gets users by slice of IDs
func (v *UserDAO) BatchGet(ctx context.Context, ids []string) ([]*dto.User, error) {
	var users []*dto.User
	for _, id := range ids {
		user, err := v.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// Query queries users by sort, range, filter
func (v *UserDAO) Query(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.User, error) {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Users)

	var cursor *mongo.Cursor
	var err error
	var count int64
	var users []*dto.User

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
			return 1, []*dto.User{d}, nil
		}

		// else: do query filter
		if filter.Item == "q" {
			query := bson.M{
				"$or": bson.A{
					bson.M{
						constants.Name: bson.M{
							"$regex":   fmt.Sprintf("%s.*", filter.Value),
							"$options": "i",
						},
					},
					bson.M{
						constants.Email: bson.M{
							"$regex":   fmt.Sprintf("%s.*", filter.Value),
							"$options": "i",
						},
					},
					bson.M{
						constants.PhoneNumber: bson.M{
							"$regex":   fmt.Sprintf("%s.*", filter.Value),
							"$options": "i",
						},
					},
					bson.M{
						constants.IC: bson.M{
							"$regex":   fmt.Sprintf("%s.*", filter.Value),
							"$options": "i",
						},
					},
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
		user := &dto.User{}
		if err = cursor.Decode(&user); err != nil {
			return 0, nil, err
		}

		if user.Location != nil && len(user.Location.Coordinates) == 2 {
			user.Long = user.Location.Coordinates[0]
			user.Lat = user.Location.Coordinates[1]
		}

		users = append(users, user)
	}

	return count, users, nil
}

// Delete deletes user by ID
func (v *UserDAO) Delete(ctx context.Context, id string) error {
	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Users)
	if _, err := collection.DeleteOne(ctx, bson.D{{constants.ID, id}}); err != nil {
		return err
	}
	return nil
}

// BatchDelete deletes users by IDs
func (v *UserDAO) BatchDelete(ctx context.Context, ids []string) ([]string, error) {
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

// Update updates user
func (v *UserDAO) Update(ctx context.Context, user *dto.User) (*dto.User, error) {

	// create location
	location := &dto.Location{
		Type:        "Point",
		Coordinates: []float64{user.Long, user.Lat},
	}
	user.Location = location

	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Users)
	_, err := collection.UpdateOne(ctx, bson.D{{constants.ID, user.ID}}, bson.D{
		{"$set", user},
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetNearbyUsers gets users within 50 meter
func (v *UserDAO) GetNearbyUsers(ctx context.Context, user *dto.User) (int64, []*dto.User, error) {
	if user.Location == nil || len(user.Location.Coordinates) != 2 {
		return 0, nil, constants.InvalidArgumentError
	}

	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Users)

	query := bson.D{{
		"$and",
		bson.A{
			bson.D{{
				constants.IsActive,
				bson.D{{
					"$eq",
					true,
				}},
			}},
			bson.D{{
				constants.Role,
				bson.D{{
					"$eq",
					constants.User,
				}},
			}},
			bson.D{{
				constants.ID,
				bson.D{{
					"$ne",
					user.ID,
				}},
			}},
			bson.M{
				// radius set to 100 meter only
				constants.Location: bson.M{
					"$geoWithin": bson.M{
						"$centerSphere": bson.A{
							bson.A{
								user.Location.Coordinates[0], user.Location.Coordinates[1],
							}, 0.1 / 6378.1,
						},
					},
				},
			},
		},
	}}

	cursor, err := collection.Find(ctx, query)
	if err != nil {
		return 0, nil, err
	}

	var users []*dto.User
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		u := &dto.User{}
		if err = cursor.Decode(&u); err != nil {
			return 0, nil, err
		}

		if u.Location != nil && len(u.Location.Coordinates) == 2 {
			u.Long = u.Location.Coordinates[0]
			u.Lat = u.Location.Coordinates[1]
		}

		users = append(users, u)
	}

	// get nearby users (within 50 meter)
	query = bson.D{{
		"$and",
		bson.A{
			bson.D{{
				constants.IsActive,
				bson.D{{
					"$eq",
					true,
				}},
			}},
			bson.D{{
				constants.Role,
				bson.D{{
					"$eq",
					constants.User,
				}},
			}},
			bson.D{{
				constants.ID,
				bson.D{{
					"$ne",
					user.ID,
				}},
			}},
			bson.M{
				// radius set to 100 meter only
				constants.Location: bson.M{
					"$geoWithin": bson.M{
						"$centerSphere": bson.A{
							bson.A{
								user.Location.Coordinates[0], user.Location.Coordinates[1],
							}, 0.05 / 6378.1,
						},
					},
				},
			},
		},
	}}

	cursor, err = collection.Find(ctx, query)
	if err != nil {
		return 0, nil, err
	}

	var nearbyUsers []*dto.User
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		u := &dto.User{}
		if err = cursor.Decode(&u); err != nil {
			return 0, nil, err
		}

		if u.Location != nil && len(u.Location.Coordinates) == 2 {
			u.Long = u.Location.Coordinates[0]
			u.Lat = u.Location.Coordinates[1]
		}

		nearbyUsers = append(nearbyUsers, u)
	}

	// update nearby users list
	err = v.updateUsersList(ctx, user, nearbyUsers)
	if err != nil {
		return 0, nil, err
	}

	return int64(len(nearbyUsers)), users, nil
}

// updateUsersList update nearby users list for contact tracing
func (v *UserDAO) updateUsersList(ctx context.Context, user *dto.User, targetUsers []*dto.User) error {

	// call get user to get user list first
	userWithList, err := v.Get(ctx, user.ID)
	if err != nil {
		return err
	}

	collection := v.client.Database(constants.GalaSejahtera).Collection(constants.Users)

	for _, targetUser := range targetUsers {

		// clone target user into tmp
		tmp := &dto.User{
			ID:    targetUser.ID,
			Role:  targetUser.Role,
			Email: targetUser.Email,
			Lat:   targetUser.Lat,
			Long:  targetUser.Long,
			Time:  utility.TimeToMilli(utility.MalaysiaTime(time.Now())),
		}

		// if target user is in user.Users, pull the result out
		if utility.UserInUsers(userWithList.Users, tmp) {
			query := bson.M{
				constants.ID: user.ID,
			}

			update := bson.M{
				"$pull": bson.M{
					constants.Users: bson.M{
						constants.ID: tmp.ID,
					},
				},
			}

			// update user
			_, err := collection.UpdateOne(ctx, query, update)
			if err != nil {
				return err
			}
		}

		// push user into users list
		query := bson.M{
			constants.ID: user.ID,
		}

		update := bson.M{
			"$push": bson.M{
				constants.Users: tmp,
			},
		}

		// update user
		_, err := collection.UpdateOne(ctx, query, update)
		if err != nil {
			return err
		}
	}

	return nil
}
