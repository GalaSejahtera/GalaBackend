package dto

// User ...
type User struct {
	ID           string    `json:"id" bson:"id"`
	Role         string    `json:"role" bson:"role"`
	Name         string    `json:"name" bson:"name"`
	Email        string    `json:"email" bson:"email"`
	Password     string    `json:"password" bson:"password"`
	IsActive     bool      `json:"isActive" bson:"isActive"`
	LastUpdated  int64     `json:"lastUpdated" bson:"lastUpdated"`
	Lat          float64   `json:"lat" bson:"-"`
	Long         float64   `json:"long" bson:"-"`
	Consent      int64     `json:"consent" bson:"consent"`
	Location     *Location `json:"-" bson:"location"`
	AccessToken  string    `json:"accessToken" bson:"-"`
	RefreshToken string    `json:"refreshToken" bson:"-"`
	ResetToken   string    `json:"resetToken" bson:"-"`
	AccessUuid   string    `json:"accessUuid" bson:"-"`
	RefreshUuid  string    `json:"refreshUuId" bson:"-"`
	AtExpires    int64     `json:"atExpires" bson:"-"`
	RtExpires    int64     `json:"rtExpires" bson:"-"`
	ResetExpires int64     `json:"resetExpires" bson:"-"`
	Time         int64     `json:"time" bson:"time"`
	Notification bool      `json:"notification" bson:"notification"`
	Users        []*User   `json:"users" bson:"users"`
	Zones        []*Zone   `json:"zones" bson:"zones"`
	IC           string    `json:"ic" bson:"ic"`
	PhoneNumber  string    `json:"phoneNumber" bson:"phoneNumber"`
	Alert        bool      `json:"alert" bson:"alert"`
	Infected     bool      `json:"infected" bson:"infected"`
}
