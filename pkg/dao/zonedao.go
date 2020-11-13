package dao

import (
	"context"
	"fmt"
	"safeworkout/pkg/constants"
	"safeworkout/pkg/dto"
	"safeworkout/pkg/logger"
	"safeworkout/pkg/utility"
	sort2 "sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ZoneDAO ...
type ZoneDAO struct {
	client *mongo.Client
}

// InitZoneDAO ...
func InitZoneDAO(client *mongo.Client) IZoneDAO {
	return &ZoneDAO{client: client}
}

// Create creates new zone
func (v *ZoneDAO) Create(ctx context.Context, zone *dto.Zone) (*dto.Zone, error) {
	// create location
	location := &dto.Location{
		Type:        "Point",
		Coordinates: []float64{zone.Long, zone.Lat},
	}
	zone.Location = location

	// create zone
	collection := v.client.Database(constants.SafeWorkout).Collection(constants.Zones)
	if _, err := collection.InsertOne(ctx, zone); err != nil {
		return nil, err
	}
	return zone, nil
}

// Get gets zone by ID
func (v *ZoneDAO) Get(ctx context.Context, id string) (*dto.Zone, error) {
	collection := v.client.Database(constants.SafeWorkout).Collection(constants.Zones)
	zone := &dto.Zone{}
	if err := collection.FindOne(ctx, bson.D{{constants.ID, id}}).Decode(&zone); err != nil {
		return nil, err
	}

	if zone.Location != nil && len(zone.Location.Coordinates) == 2 {
		zone.Long = zone.Location.Coordinates[0]
		zone.Lat = zone.Location.Coordinates[1]
	}

	// get number of users within the zone
	users, err := v.GetUsersByZone(ctx, zone)
	if err != nil {
		logger.Log.Warn("Unable to trace users within the zone. Setting zone usersWithin to 0.")
		return zone, nil
	}
	zone.UsersWithin = int64(len(users))

	// get zone risk
	zone.Risk = utility.GetZoneRisk(zone.Radius, zone.Capacity, zone.UsersWithin)

	// get isCapacityExceeded
	if zone.Risk == constants.MaximumRisk {
		zone.IsCapacityExceeded = zone.UsersWithin > zone.Capacity
	}

	return zone, nil
}

// BatchGet gets zones by slice of IDs
func (v *ZoneDAO) BatchGet(ctx context.Context, ids []string) ([]*dto.Zone, error) {
	var zones []*dto.Zone
	for _, id := range ids {
		zone, err := v.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		zones = append(zones, zone)
	}
	return zones, nil
}

// Query queries zones by sort, range, filter
func (v *ZoneDAO) Query(ctx context.Context, sort *dto.SortData, itemsRange *dto.RangeData, filter *dto.FilterData) (int64, []*dto.Zone, error) {
	collection := v.client.Database(constants.SafeWorkout).Collection(constants.Zones)

	var cursor *mongo.Cursor
	var err error
	var count int64
	var zones []*dto.Zone

	findOptions := options.Find()
	// set range
	if itemsRange != nil {
		// special case: don't set if they are derived fields
		if sort == nil || (sort.Item != "usersWithin" && sort.Item != "isCapacityExceeded" && sort.Item != "risk") {
			findOptions.SetSkip(int64(itemsRange.From))
			findOptions.SetLimit(int64(itemsRange.To + 1 - itemsRange.From))
		}
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
			return 1, []*dto.Zone{d}, nil
		}

		// else: do query filter
		if filter.Item == "q" {
			query := bson.M{
				constants.Name: bson.M{
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
		zone := &dto.Zone{}
		if err = cursor.Decode(&zone); err != nil {
			return 0, nil, err
		}

		if zone.Location != nil && len(zone.Location.Coordinates) == 2 {
			zone.Long = zone.Location.Coordinates[0]
			zone.Lat = zone.Location.Coordinates[1]
		}

		// get number of users within the zone
		users, err := v.GetUsersByZone(ctx, zone)
		if err != nil {
			logger.Log.Warn("Unable to trace users within the zone. Setting zone usersWithin to 0.")
		} else {
			zone.UsersWithin = int64(len(users))
		}

		// get zone risk
		zone.Risk = utility.GetZoneRisk(zone.Radius, zone.Capacity, zone.UsersWithin)

		// get isCapacityExceeded
		if zone.Risk == constants.MaximumRisk {
			zone.IsCapacityExceeded = zone.UsersWithin > zone.Capacity
		}

		zones = append(zones, zone)
	}

	// special case: sort and pagination for derived fields
	if sort != nil && (sort.Item == "usersWithin" || sort.Item == "isCapacityExceeded" || sort.Item == "risk") {
		v.sortZones(zones, sort.Item, sort.Order)
		if itemsRange != nil && itemsRange.To != 0 {
			index := 0
			var rslt []*dto.Zone

			for _, z := range zones {
				if index < itemsRange.From {
					index += 1
					continue
				}
				if itemsRange.To != 0 && index > itemsRange.To {
					index += 1
					continue
				}
				index += 1
				rslt = append(rslt, z)
			}
			return count, rslt, nil
		}

	}

	return count, zones, nil
}

func (v *ZoneDAO) sortZones(zones []*dto.Zone, field string, order string) {
	switch field {
	case "usersWithin":
		sort2.Slice(zones, func(i, j int) bool {
			return zones[i].UsersWithin < zones[j].UsersWithin
		})
	case "isCapacityExceeded":
		sort2.Slice(zones, func(i, j int) bool {
			if !zones[i].IsCapacityExceeded {
				return true
			}
			return false
		})
	case "risk":
		sort2.Slice(zones, func(i, j int) bool {
			return zones[i].Risk < zones[j].Risk
		})
	default:
	}

	if order == "DESC" {
		// reverse slice
		for i, j := 0, len(zones)-1; i < j; i, j = i+1, j-1 {
			zones[i], zones[j] = zones[j], zones[i]
		}
	}
}

// Delete deletes zone by ID
func (v *ZoneDAO) Delete(ctx context.Context, id string) error {
	collection := v.client.Database(constants.SafeWorkout).Collection(constants.Zones)
	if _, err := collection.DeleteOne(ctx, bson.D{{constants.ID, id}}); err != nil {
		return err
	}
	return nil
}

// BatchDelete deletes zones by IDs
func (v *ZoneDAO) BatchDelete(ctx context.Context, ids []string) ([]string, error) {
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

// Update updates zone
func (v *ZoneDAO) Update(ctx context.Context, zone *dto.Zone) (*dto.Zone, error) {
	// update location
	location := &dto.Location{
		Type:        "Point",
		Coordinates: []float64{zone.Long, zone.Lat},
	}
	zone.Location = location

	// update zone
	collection := v.client.Database(constants.SafeWorkout).Collection(constants.Zones)
	_, err := collection.UpdateOne(ctx, bson.D{{constants.ID, zone.ID}}, bson.D{
		{"$set", zone},
	})
	if err != nil {
		return nil, err
	}
	return zone, nil
}

// GetUsersByZone gets users by zone
func (v *ZoneDAO) GetUsersByZone(ctx context.Context, zone *dto.Zone) ([]*dto.User, error) {
	if zone.Location == nil || len(zone.Location.Coordinates) != 2 {
		return []*dto.User{}, nil
	}

	collection := v.client.Database(constants.SafeWorkout).Collection(constants.Users)

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
			bson.M{
				constants.Location: bson.M{
					"$geoWithin": bson.M{
						"$centerSphere": bson.A{
							bson.A{
								zone.Location.Coordinates[0], zone.Location.Coordinates[1],
							}, zone.Radius / 6378.1,
						},
					},
				},
			},
		},
	}}

	cursor, err := collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	var users []*dto.User
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		user := &dto.User{}
		if err = cursor.Decode(&user); err != nil {
			return nil, err
		}

		if user.Location != nil && len(user.Location.Coordinates) == 2 {
			user.Long = user.Location.Coordinates[0]
			user.Lat = user.Location.Coordinates[1]
		}

		users = append(users, user)
	}

	return users, nil
}

// GetSubZones gets sub zones by zone
func (v *ZoneDAO) GetSubZones(ctx context.Context, zone *dto.Zone, user *dto.User) ([]*dto.Zone, error) {
	if zone.Location == nil || len(zone.Location.Coordinates) != 2 {
		return []*dto.Zone{}, nil
	}

	collection := v.client.Database(constants.SafeWorkout).Collection(constants.Zones)

	query := bson.D{{
		"$and",
		bson.A{
			bson.D{{
				constants.Type,
				bson.D{{
					"$eq",
					constants.SubZone,
				}},
			}},
			bson.M{
				constants.Location: bson.M{
					"$geoWithin": bson.M{
						"$centerSphere": bson.A{
							bson.A{
								zone.Location.Coordinates[0], zone.Location.Coordinates[1],
							}, zone.Radius / 6378.1,
						},
					},
				},
			},
		},
	}}

	cursor, err := collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	var zones []*dto.Zone
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		zone := &dto.Zone{}
		if err = cursor.Decode(&zone); err != nil {
			return nil, err
		}

		if zone.Location != nil && len(zone.Location.Coordinates) == 2 {
			zone.Long = zone.Location.Coordinates[0]
			zone.Lat = zone.Location.Coordinates[1]
		}

		// get number of users within the zone
		users, err := v.GetUsersByZone(ctx, zone)
		if err != nil {
			logger.Log.Warn("Unable to trace users within the zone. Setting zone usersWithin to 0.")
		} else {
			zone.UsersWithin = int64(len(users))
		}

		// get zone risk
		zone.Risk = utility.GetZoneRisk(zone.Radius, zone.Capacity, zone.UsersWithin)

		// get isCapacityExceeded
		if zone.Risk == constants.MaximumRisk {
			zone.IsCapacityExceeded = zone.UsersWithin > zone.Capacity
		}

		// if user in the sub zone, update visited list
		if utility.UserInUsers(users, user) {
			err = v.updateZonesList(ctx, user, zone)
			if err != nil {
				return nil, err
			}
			err = v.updateUsersList(ctx, zone, user)
			if err != nil {
				return nil, err
			}
		}

		zones = append(zones, zone)
	}

	return zones, nil
}

// GetByUser gets zone and sub zones given user
func (v *ZoneDAO) GetByUser(ctx context.Context, user *dto.User) (*dto.Zone, []*dto.Zone, error) {
	if user.Location == nil || len(user.Location.Coordinates) != 2 {
		return &dto.Zone{}, []*dto.Zone{}, nil
	}

	// get all zones within 100 km
	collection := v.client.Database(constants.SafeWorkout).Collection(constants.Zones)

	query := bson.D{{
		"$and",
		bson.A{
			bson.D{{
				constants.Type,
				bson.D{{
					"$eq",
					constants.Zone,
				}},
			}},
			bson.M{
				constants.Location: bson.M{
					"$nearSphere": bson.M{
						"$geometry": bson.M{
							constants.Type:        user.Location.Type,
							constants.Coordinates: user.Location.Coordinates,
						},
						"$maxDistance": 100 * 1000,
					},
				},
			},
		},
	}}

	cursor, err := collection.Find(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	var zones []*dto.Zone
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		zone := &dto.Zone{}
		if err = cursor.Decode(&zone); err != nil {
			return nil, nil, err
		}
		zones = append(zones, zone)
	}

	// for each zone from nearest to furthest, try to get zone that included the target user
	for _, z := range zones {
		u, err := v.GetUsersByZone(ctx, z)
		if err != nil {
			return nil, nil, err
		}
		// if user in users
		if utility.UserInUsers(u, user) {
			if z.Location != nil && len(z.Location.Coordinates) == 2 {
				z.Long = z.Location.Coordinates[0]
				z.Lat = z.Location.Coordinates[1]
			}

			// set number of users
			z.UsersWithin = int64(len(u))

			// get zone risk
			z.Risk = utility.GetZoneRisk(z.Radius, z.Capacity, z.UsersWithin)

			// get isCapacityExceeded
			if z.Risk == constants.MaximumRisk {
				z.IsCapacityExceeded = z.UsersWithin > z.Capacity
			}

			// update user visited zone list for main zone
			err = v.updateZonesList(ctx, user, z)
			if err != nil {
				return nil, nil, err
			}
			err = v.updateUsersList(ctx, z, user)
			if err != nil {
				return nil, nil, err
			}

			// get sub zones
			subZones, err := v.GetSubZones(ctx, z, user)
			if err != nil {
				return nil, nil, err
			}

			return z, subZones, nil
		}
	}

	return &dto.Zone{}, []*dto.Zone{}, nil
}

// QueryRecentUsersByZoneID queries past 14 days users by zoneID
func (v *ZoneDAO) QueryRecentUsersByZoneID(ctx context.Context, zoneID string) ([]*dto.User, error) {
	zone, err := v.Get(ctx, zoneID)
	if err != nil {
		return nil, err
	}

	// set timestamp to 14 days ago, 12 am
	now := utility.MalaysiaTime(time.Now())
	daySelected, err := utility.DateStringToTime(utility.TimeToDateString(now.Add(time.Duration(-13) * 24 * time.Hour)))
	t := utility.TimeToMilli(daySelected)
	if err != nil {
		return nil, err
	}

	var users []*dto.User

	for _, user := range zone.Users {
		if user.Time >= t {
			users = append(users, user)
		}
	}

	sort2.Slice(users, func(i, j int) bool {
		return users[i].Time > users[j].Time
	})

	return users, nil
}

// updateZonesList update visited zones list for contact tracing
func (v *ZoneDAO) updateZonesList(ctx context.Context, user *dto.User, targetZone *dto.Zone) error {
	collection := v.client.Database(constants.SafeWorkout).Collection(constants.Users)

	// clone target zone into tmp
	tmp := &dto.Zone{
		ID:                 targetZone.ID,
		Name:               targetZone.Name,
		Lat:                targetZone.Lat,
		Long:               targetZone.Long,
		Type:               targetZone.Type,
		Capacity:           targetZone.Capacity,
		Radius:             targetZone.Radius,
		UsersWithin:        targetZone.UsersWithin,
		IsCapacityExceeded: targetZone.IsCapacityExceeded,
		Time:               utility.TimeToMilli(utility.MalaysiaTime(time.Now())),
		Risk:               utility.GetZoneRisk(targetZone.Radius, targetZone.Capacity, targetZone.UsersWithin),
	}

	// if target zone is in user.Zones, pull the result out
	if utility.ZoneInZones(user.Zones, tmp) {

		query := bson.M{
			constants.ID: user.ID,
		}

		update := bson.M{
			"$pull": bson.M{
				constants.Zones: bson.M{
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

	// push zone into zones list
	query := bson.M{
		constants.ID: user.ID,
	}

	update := bson.M{
		"$push": bson.M{
			constants.Zones: tmp,
		},
	}

	// update user
	_, err := collection.UpdateOne(ctx, query, update)
	if err != nil {
		return err
	}

	return nil
}

// updateUsersList update visited users list for contact tracing
func (v *ZoneDAO) updateUsersList(ctx context.Context, zone *dto.Zone, targetUser *dto.User) error {
	// call get zone to get user list first
	zoneWithList, err := v.Get(ctx, zone.ID)
	if err != nil {
		return err
	}

	collection := v.client.Database(constants.SafeWorkout).Collection(constants.Zones)

	// clone target user into tmp
	tmp := &dto.User{
		ID:          targetUser.ID,
		Role:        targetUser.Role,
		Name:        targetUser.Name,
		Email:       targetUser.Email,
		Lat:         targetUser.Lat,
		Long:        targetUser.Long,
		Time:        utility.TimeToMilli(utility.MalaysiaTime(time.Now())),
		IC:          targetUser.IC,
		PhoneNumber: targetUser.PhoneNumber,
		Infected:    targetUser.Infected,
	}

	// if target user is in zone.Users, pull the result out
	if utility.UserInUsers(zoneWithList.Users, tmp) {

		query := bson.M{
			constants.ID: zone.ID,
		}

		update := bson.M{
			"$pull": bson.M{
				constants.Users: bson.M{
					constants.ID: tmp.ID,
				},
			},
		}

		// update zone
		_, err := collection.UpdateOne(ctx, query, update)
		if err != nil {
			return err
		}
	}

	// push user into users list
	query := bson.M{
		constants.ID: zone.ID,
	}

	update := bson.M{
		"$push": bson.M{
			constants.Users: tmp,
		},
	}

	// update user
	_, err = collection.UpdateOne(ctx, query, update)
	if err != nil {
		return err
	}

	return nil
}
