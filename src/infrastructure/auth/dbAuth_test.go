package auth

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSessionManager_NewSessionManager(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	c, err := redis.Dial("tcp", s.Addr())
	sessionManager := NewSessionManager(c)

	require.EqualValues(t, c, sessionManager.RedisConn)
}

func TestSessionManager_CreateAuthSuccess(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	c, err := redis.Dial("tcp", s.Addr())
	sessionManager := NewSessionManager(c)

	cookie, err := sessionManager.CreateAuth(1)
	require.NoError(t, err)
	require.NotEmpty(t, cookie)

	if got, err := s.Get("sessions:" + cookie.Value); err != nil || got != "1" {
		t.Error("'foo' has the wrong value")
	}
}

func TestNewSessionManager_CheckAuthSuccess(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	c, err := redis.Dial("tcp", s.Addr())
	sessionManager := NewSessionManager(c)

	s.Set("sessions:session_value", "1")
	id, err := sessionManager.CheckAuth("session_value")
	require.NoError(t, err)
	require.EqualValues(t, 1, id)
}

func TestNewSessionManager_CheckAuthFail(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	c, err := redis.Dial("tcp", s.Addr())
	sessionManager := NewSessionManager(c)

	id, err := sessionManager.CheckAuth("session_value")
	require.EqualValues(t, "no session", err.Error())
	require.EqualValues(t, 0, id)
}

func TestNewSessionManager_DeleteAuth (t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	c, err := redis.Dial("tcp", s.Addr())
	sessionManager := NewSessionManager(c)

	s.Set("sessions:session_value", "1")
	err = sessionManager.DeleteAuth("session_value")
	require.NoError(t, err)

	exists := s.Exists("session_value")
	require.EqualValues(t, false, exists)
}


