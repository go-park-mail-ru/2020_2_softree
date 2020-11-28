package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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

func MakeShieldedCookie() (string, error) {
	hash := sha256.New()

	salt := make([]byte, 8)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	if _, err := hash.Write(salt); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
