package security

import (
	"golang.org/x/crypto/bcrypt"
)

type Utils struct {
	defaultCost int
}

func CreateNewSecurityUtils() *Utils {
	return &Utils{defaultCost: 10}
}

func (u *Utils) MakeShieldedPassword(stringToHash string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(stringToHash), u.defaultCost)

	return string(pass), err
}
