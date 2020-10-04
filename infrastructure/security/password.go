package security

import (
	"crypto/md5"
	"encoding/hex"
)

func MakeDoubleHash(stringToHash string) string {
	hash := md5.New()
	hash.Write([]byte(stringToHash))
	hash.Write([]byte(hex.EncodeToString(hash.Sum(nil))))

	return hex.EncodeToString(hash.Sum(nil))
}
