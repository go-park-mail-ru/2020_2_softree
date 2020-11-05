package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Session struct {
	id    uint64
	value string
}

type SessionManager struct {
	redisConn redis.Conn
}

func NewSessionManager(conn redis.Conn) *SessionManager {
	return &SessionManager{
		redisConn: conn,
	}
}

func (sm *SessionManager) CreateAuth(id uint64) (cookie http.Cookie, err error) {
	if cookie, err = sm.CreateCookie(); err != nil {
		return http.Cookie{}, err
	}

	mkey := "sessions:" + cookie.Value
	result, err := redis.String(sm.redisConn.Do("SET", mkey, id, "EX", 60*60*24)) // Expires in 24 hours
	if err != nil {
		return http.Cookie{}, err
	}
	if result != "OK" {
		return http.Cookie{}, fmt.Errorf("result not OK")
	}

	return cookie, nil
}

func (sm *SessionManager) CheckAuth(sessionValue string) (uint64, error) {
	mkey := "sessions:" + sessionValue
	data, err := redis.Bytes(sm.redisConn.Do("GET", mkey))
	if err != nil {
		return 0, errors.New("redis error during checking session")
	}
	if data == "nil" {
		return 0, errors.New("no session")
	}

	return uint64(data), nil
}

func (sm *SessionManager) DeleteAuth(sessionValue string) error {
	mkey := "sessions:" + sessionValue
	_, err := redis.Int(sm.redisConn.Do("DEL", mkey))
	if err != nil {
		return errors.New("redis error during session delete")
	}

	return nil
}

func (sm *SessionManager) CreateCookie() (http.Cookie, error) {
	hash, err := security.MakeShieldedCookie()
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
