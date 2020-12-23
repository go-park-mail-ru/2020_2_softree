package authorization

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/infrastructure/mock"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	id       = int64(1)
	email    = "hound@psina.ru"
	password = "str"
	avatar   = "base64"

	name  = "session_id"
	value = "value"
)

func TestAuth_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/auth"
	body := strings.NewReader(fmt.Sprintf("{\"email\": %s, \"password\": %s}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), entity.UserIdKey, int64(id))
	req = req.Clone(ctx)
	testAuth, ctrl := createAuthSuccess(t, req.Context())
	defer ctrl.Finish()

	req.AddCookie(&http.Cookie{Name: name, Value: value})

	testAuth.Authenticate(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-type"))
	require.NotEmpty(t, w.Body)
}

func TestAuth_FailUnauthorized(t *testing.T) {
	url := "http://127.0.0.1:8000/auth"
	body := strings.NewReader(fmt.Sprintf("{\"email\": %s, \"password\": %s}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createAuthFailUnauthorized(t)
	defer ctrl.Finish()

	emptyFunc := testAuth.Auth(empty)
	emptyFunc(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}

func TestAuthCheck_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/auth"
	body := strings.NewReader(fmt.Sprintf("{\"email\": %s, \"password\": %s}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createAuthCheckSuccess(t, req.Context())
	defer ctrl.Finish()

	cookie := http.Cookie{
		Name:  name,
		Value: value,
	}
	req.AddCookie(&cookie)

	emptyFunc := testAuth.Auth(empty)
	emptyFunc(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestAuth_FailNoSession(t *testing.T) {
	url := "http://127.0.0.1:8000/auth"
	body := strings.NewReader(fmt.Sprintf("{\"email\": %s, \"password\": %s}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	ctx := req.Context()
	testAuth, ctrl := createAuthFailSession(t, ctx)
	defer ctrl.Finish()

	cookie := http.Cookie{
		Name:  name,
		Value: value,
	}
	req.AddCookie(&cookie)

	auth := testAuth.Auth(testAuth.Authenticate)
	auth(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func createAuthSuccess(t *testing.T, ctx context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := mock.NewMockAuthLogic(ctrl)
	mockAuth.EXPECT().
		Authenticate(ctx, id).
		Return(entity.Description{}, createExpectedUser(), nil)

	mockUser := mock.NewMockProfileLogic(ctrl)

	return NewAuthentication(mockUser, mockAuth), ctrl
}

func createAuthFailUnauthorized(t *testing.T) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mock.NewMockProfileLogic(ctrl)
	mockAuth := mock.NewMockAuthLogic(ctrl)

	return NewAuthentication(mockUser, mockAuth), ctrl
}

func createAuthCheckSuccess(t *testing.T, ctx context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mock.NewMockProfileLogic(ctrl)
	mockAuth := mock.NewMockAuthLogic(ctrl)

	mockAuth.EXPECT().
		Auth(ctx, &http.Cookie{Name: name, Value: value}).
		Return(entity.Description{}, id, nil)

	return NewAuthentication(mockUser, mockAuth), ctrl
}

func createAuthFailSession(t *testing.T, ctx context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mock.NewMockProfileLogic(ctrl)

	mockAuth := mock.NewMockAuthLogic(ctrl)
	mockAuth.EXPECT().
		Auth(ctx, &http.Cookie{Name: name, Value: value}).
		Return(entity.Description{Status: http.StatusBadRequest}, int64(0), errors.New("no session"))

	return NewAuthentication(mockUser, mockAuth), ctrl
}

func empty(w http.ResponseWriter, r *http.Request) {
	UserId := r.Context().Value(entity.UserIdKey).(int64)

	if UserId == id {
		w.WriteHeader(http.StatusOK)
	}
	w.WriteHeader(http.StatusBadRequest)
}

func createExpectedUser() entity.PublicUser {
	return entity.PublicUser{Id: id, Email: email, Avatar: avatar}
}
