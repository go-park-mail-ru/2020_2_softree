package authorization

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	authMock "server/src/authorization/pkg/infrastructure/mock"
	session "server/src/authorization/pkg/session/gen"
	profileMock "server/src/profile/pkg/infrastructure/mock"
	profile "server/src/profile/pkg/profile/gen"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	id = 1
	email = "hound@psina.ru"
	avatar = "base64"

	name = "session_id"
	value = "value"
)

func TestAuth_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/auth"
	body := strings.NewReader(`{"email": , "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	ctx := req.Context()
	testAuth, ctrl := createAuthSuccess(t, &ctx)
	defer ctrl.Finish()

	req.AddCookie(&http.Cookie{Name: name, Value: value})

	auth := testAuth.Auth(testAuth.Authenticate)
	auth(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-type"))
	require.NotEmpty(t, w.Body)
}

func TestAuth_FailUnauthorized(t *testing.T) {
	url := "http://127.0.0.1:8000/auth"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createAuthFailUnauthorized(t)
	defer ctrl.Finish()

	auth := testAuth.Auth(testAuth.Authenticate)
	auth(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}

func TestAuth_FailNoSession(t *testing.T) {
	url := "http://127.0.0.1:8000/auth"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createAuthFailSession(t)
	defer ctrl.Finish()

	cookie := http.Cookie{
		Name:  "session_id",
		Value: "value",
	}
	req.AddCookie(&cookie)

	auth := testAuth.Auth(testAuth.Authenticate)
	auth(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestAuth_FailNoUser(t *testing.T) {
	url := "http://127.0.0.1:8000/auth"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createAuthFailUser(t)
	defer ctrl.Finish()

	cookie := http.Cookie{
		Name:  "session_id",
		Value: "value",
	}
	req.AddCookie(&cookie)

	auth := testAuth.Auth(testAuth.Authenticate)
	auth(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func createAuthSuccess(t *testing.T, ctx *context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().GetUserById(ctx, &profile.UserID{Id: id}).Return(createExpectedUser(), nil)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().Check(ctx, &session.SessionID{SessionId: value}).Return(&session.UserID{Id: id}, nil)

	return NewAuthenticate(mockUser, mockAuth), ctrl
}

func createAuthFailUnauthorized(t *testing.T) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)
	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)

	return NewAuthenticate(mockUser, mockAuth), ctrl
}

func createAuthFailSession(t *testing.T, ctx *context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().Check(ctx, &session.SessionID{SessionId: value}).Return(int64(0), errors.New("no session"))

	return NewAuthenticate(mockUser, mockAuth), ctrl
}

func createAuthFailUser(t *testing.T, ctx *context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().GetUserById(ctx, &profile.UserID{Id: id}).Return(nil, errors.New("no user in database"))

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().Check(ctx, &session.SessionID{SessionId: value}).Return(&session.UserID{Id: id}, nil)

	return NewAuthenticate(mockUser, mockAuth), ctrl
}

func createExpectedUser() *profile.User {
	return &profile.User{ID: id, Email: email, Avatar: avatar}
}
