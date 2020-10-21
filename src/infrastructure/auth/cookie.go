package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"net/http"
	"server/src/infrastructure/config"
	"strconv"
	"time"
)

type CookieInterface interface {
	CreateCookie(uint64) (*CookieDetails, error)
	ExtractData(*http.Request) (*AccessDetails, error)
}

type Cookie struct{}

func (c *Cookie) CreateCookie() (*CookieDetails, error) {
	hash, err := makeCookieHash()
	if err != nil {
		return &CookieDetails{}, err
	}

	cd := MakeCookieDetailsFromCookie(http.Cookie{
		Name:     "session_id",
		Value:    hash,
		Expires:  time.Now().Add(24 * time.Hour),
		Domain:   config.GlobalServerConfig.Domain,
		Secure:   config.GlobalServerConfig.Secure,
		HttpOnly: true,
		Path:     "/",
	})

	return &cd, nil
}

func (c *Cookie) ExtractData(r *http.Request) (ad *AccessDetails, err error) {
	if err := json.NewDecoder(r.Body).Decode(&ad); err != nil {
		return &AccessDetails{}, err
	}
	return ad, nil
}

func makeCookieHash() (string, error) {
	hash := sha256.New()

	if _, err := hash.Write([]byte(strconv.Itoa(rand.Int()))); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
