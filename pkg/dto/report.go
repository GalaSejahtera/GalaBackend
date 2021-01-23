package dto

// Report ...
type Report struct {
	ID         string `json:"id" bson:"id"`
	UserID     string `json:"userId" bson:"userId"`
	CreatedAt  int64  `json:"createdAt" bson:"createdAt"`
	HasSymptom bool   `json:"hasSymptom" bson:"hasSymptom"`
	Results    []bool `json:"results" bson:"results"`
}
