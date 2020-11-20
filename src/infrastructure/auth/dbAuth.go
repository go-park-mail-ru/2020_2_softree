package auth

import (
	"errors"
	"net/http"
	"server/src/infrastructure/security"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

const day = 60 * 60 * 24

type Session struct {
	id    uint64
	value string
}

type SessionManager struct {
	RedisConn redis.Conn
}

func NewSessionManager(conn redis.Conn) *SessionManager {
	return &SessionManager{
		RedisConn: conn,
	}
}

func (sm *SessionManager) CreateAuth(id uint64) (cookie http.Cookie, err error) {
	if cookie, err = sm.CreateCookie(); err != nil {
		return http.Cookie{}, err
	}

	mkey := "sessions:" + cookie.Value
	result, err := redis.String(sm.RedisConn.Do("SET", mkey, id, "EX", day)) // Expires in 24 hours
	if err != nil {
		return http.Cookie{}, err
	}
	if result != "OK" {
		return http.Cookie{}, errors.New("result not OK")
	}

	return cookie, nil
}

func (sm *SessionManager) CheckAuth(sessionValue string) (uint64, error) {
	mkey := "sessions:" + sessionValue
	data, err := redis.Bytes(sm.RedisConn.Do("GET", mkey))
	if err == redis.ErrNil {
		return 0, errors.New("no session")
	} else if err != nil {
		return 0, err
	}

	strRes := string(data)
	uintRes, parseErr := strconv.ParseUint(strRes, 10, 64)
	if parseErr != nil {
		return 0, errors.New("internal server error")
	}

	return uintRes, nil
}

func (sm *SessionManager) DeleteAuth(sessionValue string) error {
	mkey := "sessions:" + sessionValue
	_, err := redis.Int(sm.RedisConn.Do("DEL", mkey))
	if err != nil {
		return err
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
		Domain:   viper.GetString("server.domain"),
		Secure:   viper.GetBool("server.secure"),
		HttpOnly: true,
		Path:     "/",
	}, nil
}
