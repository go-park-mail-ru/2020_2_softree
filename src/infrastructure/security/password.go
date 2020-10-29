package security

import (
	"crypto/md5"
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
)

var defaultCost = 10

func MakeShieldedPassword(stringToHash string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(stringToHash), defaultCost)

	return string(pass), err
}

func MakeShieldedCookie() (string, error) {
	salt := make([]byte, 8)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	return string(md5.New().Sum(salt)), nil
}
