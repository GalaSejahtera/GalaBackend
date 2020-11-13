package dto

// Zone ...
type Zone struct {
	ID                 string    `json:"id" bson:"id"`
	Name               string    `json:"name" bson:"name"`
	Lat                float64   `json:"lat" bson:"-"`
	Long               float64   `json:"long" bson:"-"`
	Type               int64     `json:"type" bson:"type"`
	Capacity           int64     `json:"capacity" bson:"capacity"`
	Radius             float64   `json:"radius" bson:"radius"`
	Location           *Location `json:"-" bson:"location"`
	UsersWithin        int64     `json:"usersWithin" bson:"-"`
	IsCapacityExceeded bool      `json:"isCapacityExceeded" bson:"-"`
	Time               int64     `json:"time" bson:"time"`
	Risk               int64     `json:"risk" bson:"risk"`
	Users              []*User   `json:"users" bson:"users"`
}
