package dto

import "time"

// AuthObject ...
type AuthObject struct {
	Token       string    `json:"token" bson:"token"`
	UserId      string    `json:"userId" bson:"userId"`
	DisplayName string    `json:"displayName" bson:"displayName"`
	Type        string    `json:"-" bson:"type"`
	TTL         time.Time `json:"ttl" bson:"ttl"`
}
