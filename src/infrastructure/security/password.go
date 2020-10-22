package security

import (
	"crypto/sha256"
	"encoding/hex"
)

func MakeShieldedHash(stringToHash string) (string, error) {
	hash := sha256.New()
	salt := "someSalt"

	stringPlusSalt := stringToHash + salt

	if _, err := hash.Write([]byte(stringPlusSalt)); err != nil {
		return "", err
	}

	if _, err := hash.Write([]byte(hex.EncodeToString(hash.Sum(nil)))); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
