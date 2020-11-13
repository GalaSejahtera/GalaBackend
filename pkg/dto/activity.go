package dto

import "time"

// Activity ...
type Activity struct {
	ID     string    `json:"id" bson:"id"`
	ZoneID string    `json:"zoneId" bson:"zoneId"`
	UserID string    `json:"userId" bson:"userId"`
	Time   int64     `json:"time" bson:"time"`
	TTL    time.Time `json:"ttl" bson:"ttl"`
}
