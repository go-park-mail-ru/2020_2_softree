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
	CreateCookie() (http.Cookie, error)
	ExtractData(*http.Request) (*AccessDetails, error)
}

type Token struct {
	token string
}

func NewToken(token string) *Token {
	return &Token{token: token}
}

func (t *Token) ExtractData(r *http.Request) (ad *AccessDetails, err error) {
	if err := json.NewDecoder(r.Body).Decode(&ad); err != nil {
		return &AccessDetails{}, err
	}
	return ad, nil
}

func (t *Token) CreateCookie() (http.Cookie, error) {
	hash, err := makeCookieHash()
	if err != nil {
		return http.Cookie{}, err
	}
	return http.Cookie{
		Name:     "session_id",
		Value:    hash,
		Expires:  time.Now().Add(24 * time.Hour),
		Domain:   config.GlobalServerConfig.Domain,
		Secure:   config.GlobalServerConfig.Secure,
		HttpOnly: true,
		Path:     "/",
	}, nil
}

func CreateCookie() (http.Cookie, error) {
	hash, err := makeCookieHash()
	if err != nil {
		return http.Cookie{}, err
	}
	return http.Cookie{
		Name:     "session_id",
		Value:    hash,
		Expires:  time.Now().Add(24 * time.Hour),
		Domain:   config.GlobalServerConfig.Domain,
		Secure:   config.GlobalServerConfig.Secure,
		HttpOnly: true,
		Path:     "/",
	}, nil
}

func makeCookieHash() (string, error) {
	hash := sha256.New()

	if _, err := hash.Write([]byte(strconv.Itoa(rand.Int()))); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
