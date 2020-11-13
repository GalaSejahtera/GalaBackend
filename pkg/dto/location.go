package dto

// Location ...
type Location struct {
	Type        string    `json:"-" bson:"type"`
	Coordinates []float64 `json:"-" bson:"coordinates"`
}
