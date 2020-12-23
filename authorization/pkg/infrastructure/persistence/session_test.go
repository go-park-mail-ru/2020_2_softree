package persistence_test

import (
	"context"
	"github.com/sirupsen/logrus"
	database "server/authorization/pkg/infrastructure/persistence"
	session "server/authorization/pkg/session/gen"
	"strconv"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/require"
)

const (
	sessionId = "some_value"
	userId    = 1
)

func TestCreate_Success(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	testPool := &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", s.Addr())
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "main",
					"action":   "connect to redis",
				}).Fatalln(err)
			}
			return conn, err
		},
	}

	sessionManager := database.NewSessionManager(testPool)

	ctx := context.Background()
	sess, err := sessionManager.Create(ctx, &session.UserID{Id: userId})
	require.NoError(t, err)
	require.NotEmpty(t, sess)

	if got, err := s.Get("sessions:" + sess.SessionId); err != nil || got != strconv.Itoa(userId) {
		t.Error("'foo' has the wrong value")
	}
}

func TestCheck_Success(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	testPool := &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", s.Addr())
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "main",
					"action":   "connect to redis",
				}).Fatalln(err)
			}
			return conn, err
		},
	}

	sessionManager := database.NewSessionManager(testPool)

	ctx := context.Background()
	require.NoError(t, s.Set("sessions:"+sessionId, strconv.Itoa(userId)))

	id, err := sessionManager.Check(ctx, &session.SessionID{SessionId: sessionId})
	require.NoError(t, err)

	require.NoError(t, err)
	require.EqualValues(t, userId, id.Id)
}

func TestCheck_Fail(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	testPool := &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", s.Addr())
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "main",
					"action":   "connect to redis",
				}).Fatalln(err)
			}
			return conn, err
		},
	}

	sessionManager := database.NewSessionManager(testPool)

	ctx := context.Background()
	_, err = sessionManager.Check(ctx, &session.SessionID{SessionId: sessionId})

	require.EqualValues(t, "reply is nil", err.Error())
}

func TestDelete_Success(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	testPool := &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", s.Addr())
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "main",
					"action":   "connect to redis",
				}).Fatalln(err)
			}
			return conn, err
		},
	}

	sessionManager := database.NewSessionManager(testPool)

	ctx := context.Background()
	require.NoError(t, s.Set("sessions:"+sessionId, strconv.Itoa(userId)))
	_, err = sessionManager.Delete(ctx, &session.SessionID{SessionId: sessionId})
	require.NoError(t, err)

	exists := s.Exists(sessionId)
	require.EqualValues(t, false, exists)
}
