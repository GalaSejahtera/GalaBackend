package dto

// Faq ...
type Faq struct {
	ID    string `json:"id" bson:"id"`
	Title string `json:"title" bson:"title"`
	Desc  string `json:"desc" bson:"desc"`
}
