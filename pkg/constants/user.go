package constants

const (
	AccessTokenTTLMinutes     = 15
	RefreshTokenTTLDays       = 365 // one year only expire the user
	PasswordResetTokenTTLDays = 1
)

const (
	SuperUser = "superuser"
	Admin     = "admin"
	User      = "user"
)

var SuperUserOnly = []string{SuperUser}
var SuperUserAndAdmin = []string{SuperUser, Admin}
var AllCanAccess = []string{SuperUser, Admin, User}
