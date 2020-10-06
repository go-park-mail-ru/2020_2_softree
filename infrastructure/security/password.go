package security

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

func MakeShieldedHash(stringToHash string) string {
	hash := sha256.New()
	salt := "someSalt"

	stringPlusSalt := stringToHash + salt

	hash.Write([]byte(stringPlusSalt))
	hash.Write([]byte(hex.EncodeToString(hash.Sum(nil))))

	return hex.EncodeToString(hash.Sum(nil))
}

func MakeCookieHash() string {
	rand.Seed(time.Now().UnixNano())
	hash := sha256.New()

	hash.Write([]byte(strconv.Itoa(rand.Int())))
	return hex.EncodeToString(hash.Sum(nil))
}
