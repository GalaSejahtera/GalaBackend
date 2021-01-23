package dto

// Covid ...
type Covid struct {
	ID              string `json:"id" bson:"id"`
	Title           string `json:"title" bson:"title"`
	SID             int64  `json:"sid" bson:"sid"`
	ImageFeatSingle string `json:"image_feat_single" bson:"image_feat_single"`
	Summary         string `json:"summary" bson:"summary"`
	DatePub2        string `json:"date_pub2" bson:"date_pub2"`
	Content         string `json:"content" bson:"content"`
	NewsURL         string `json:"newsUrl" bson:"newsUrl"`
}

// Container ...
type Container struct {
	Stories []*Covid `json:"stories" bson:"stories"`
}

// District ...
type District struct {
	Name  string `json:"name" bson:"name"`
	Total int64  `json:"total" bson:"total"`
}

// State ...
type State struct {
	Name      string      `json:"name" bson:"name"`
	Total     int64       `json:"total" bson:"total"`
	Districts []*District `json:"districts" bson:"districts"`
}

// Daily ...
type Daily struct {
	States      []*State `json:"states" bson:"states"`
	LastUpdated int64    `json:"last_updated" bson:"last_updated"`
}

// General ...
type General struct {
	TotalConfirmed int64 `json:"totalConfirmed" bson:"totalConfirmed"`
	ActiveCases    int64 `json:"activeCases" bson:"activeCases"`
}
