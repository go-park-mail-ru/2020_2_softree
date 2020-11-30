package persistence

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	session "server/authorization/pkg/session/gen"
	"strconv"
)

const day = 60 * 60 * 24

type SessionManager struct {
	RedisConn redis.Conn
}

func NewSessionManager(conn redis.Conn) *SessionManager {
	return &SessionManager{
		RedisConn: conn,
	}
}

func (sm *SessionManager) Create(ctx context.Context, in *session.UserID) (*session.Session, error) {
	hash, err := makeSessionValue()
	if err != nil {
		return nil, err
	}

	key := "sessions:" + hash
	result, err := redis.String(sm.RedisConn.Do("SET", key, in.Id, "EX", day)) // Expires in 24 hours

	if err != nil {
		return nil, err
	}
	if result != "OK" {
		return nil, errors.New("result not OK")
	}

	return &session.Session{SessionId: hash, Id: in.Id}, nil
}

func (sm *SessionManager) Check(ctx context.Context, in *session.SessionID) (*session.UserID, error) {
	key := "sessions:" + in.SessionId
	data, err := redis.Bytes(sm.RedisConn.Do("GET", key))

	if err == redis.ErrNil {
		return nil, errors.New("no session")
	} else if err != nil {
		return nil, err
	}

	var id int64
	if id, err = strconv.ParseInt(string(data), 10, 64); err != nil {
		return nil, err
	}

	return &session.UserID{Id: id}, nil
}

func (sm *SessionManager) Delete(ctx context.Context, in *session.SessionID) (*session.Empty, error) {
	key := "sessions:" + in.SessionId
	_, err := redis.Int(sm.RedisConn.Do("DEL", key))
	if err != nil {
		return &session.Empty{}, err
	}

	return &session.Empty{}, nil
}

func makeSessionValue() (string, error) {
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
