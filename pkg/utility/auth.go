package utility

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword is used to hash password using bcrypt
func HashPassword(password string) (string, error) {
	passwordByte := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
