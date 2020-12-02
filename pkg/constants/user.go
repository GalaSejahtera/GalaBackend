package constants

const (
	AccessTokenTTLMinutes     = 15
	RefreshTokenTTLDays       = 365 // one year only expire the user
	PasswordResetTokenTTLDays = 1
)

const (
	User = "user"
)

var AllCanAccess = []string{User}
