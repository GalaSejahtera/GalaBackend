package constants

// Database
const (
	GalaSejahtera = "galasejahtera"
)

// Collections
const (
	Users      = "users"
	AuthTokens = "authtokens"
	Zones      = "zones"
	Activities = "activities"
	Faqs       = "faqs"
)

// Fields
const (
	// Token
	Token       = "token"
	TTL         = "ttl"
	Authorized  = "authorized"
	AccessUuid  = "accessUuid"
	UserId      = "userId"
	Exp         = "exp"
	RefreshUuid = "refreshUuid"
	Access      = "access"
	Refresh     = "refresh"
	Reset       = "reset"

	// Common
	ID = "id"

	// Users
	DisplayName = "displayName"
	Email       = "email"
	IsActive    = "isActive"
	Role        = "role"
	Disabled    = "disabled"
	Password    = "password"
	IC          = "ic"
	PhoneNumber = "phoneNumber"
	Alert       = "alert"
	Infected    = "infected"

	// Zones
	Name     = "name"
	Location = "location"

	// Location
	Coordinates = "coordinates"
	Type        = "type"

	// Activity
	Time   = "time"
	ZoneId = "zoneId"

	// Faq
	Title = "title"
	Desc  = "desc"
)

// Keywords
const (
	ASC  = "ASC"
	DESC = "DESC"
)
