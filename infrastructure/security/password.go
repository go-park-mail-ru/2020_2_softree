package security

import (
	"crypto/sha256"
	"encoding/hex"
)

func MakeShieldedHash(stringToHash string) string {
	hash := sha256.New()
	salt := "someSalt"

	stringPlusSalt := stringToHash + salt

	hash.Write([]byte(stringPlusSalt))
	hash.Write([]byte(hex.EncodeToString(hash.Sum(nil))))

	return hex.EncodeToString(hash.Sum(nil))
}
