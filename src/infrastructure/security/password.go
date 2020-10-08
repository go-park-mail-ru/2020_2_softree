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

	_, err := hash.Write([]byte(stringPlusSalt))
	if err != nil {
		panic(err)
	}

	_, err = hash.Write([]byte(hex.EncodeToString(hash.Sum(nil))))
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(hash.Sum(nil))
}

func makeCookieHash() string {
	rand.Seed(time.Now().UnixNano())
	hash := sha256.New()

	_, err := hash.Write([]byte(strconv.Itoa(rand.Int())))
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}
