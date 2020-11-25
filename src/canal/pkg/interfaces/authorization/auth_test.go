package authorization

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	authMock "server/src/authorization/pkg/infrastructure/mock"
	session "server/src/authorization/pkg/session/gen"
	"server/src/canal/pkg/infrastructure/mock"
	profileMock "server/src/profile/pkg/infrastructure/mock"
	profile "server/src/profile/pkg/profile/gen"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	id       = 1
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

	ctx := context.WithValue(req.Context(), "id", int64(id))
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

func TestAuth_FailNoUser(t *testing.T) {
	url := "http://127.0.0.1:8000/auth"
	body := strings.NewReader(fmt.Sprintf("{\"email\": %s, \"password\": %s}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), "id", int64(id))
	req = req.Clone(ctx)
	testAuth, ctrl := createAuthFailUser(t, ctx)
	defer ctrl.Finish()

	cookie := http.Cookie{
		Name:  name,
		Value: value,
	}
	req.AddCookie(&cookie)

	testAuth.Authenticate(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func createAuthSuccess(t *testing.T, ctx context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		GetUserById(ctx, &profile.UserID{Id: id}).
		Return(createExpectedUser(), nil)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)

	return NewAuthenticate(mockUser, mockAuth, mock.NewSecurityMock(ctrl)), ctrl
}

func createAuthFailUnauthorized(t *testing.T) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)
	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)

	return NewAuthenticate(mockUser, mockAuth, mock.NewSecurityMock(ctrl)), ctrl
}

func createAuthCheckSuccess(t *testing.T, ctx context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)
	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().
		Check(ctx, &session.SessionID{SessionId: value}).
		Return(&session.UserID{Id: id}, nil)

	return NewAuthenticate(mockUser, mockAuth, mock.NewSecurityMock(ctrl)), ctrl
}

func createAuthFailSession(t *testing.T, ctx context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().
		Check(ctx, &session.SessionID{SessionId: value}).
		Return(nil, errors.New("no session"))

	return NewAuthenticate(mockUser, mockAuth, mock.NewSecurityMock(ctrl)), ctrl
}

func createAuthFailUser(t *testing.T, ctx context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		GetUserById(ctx, &profile.UserID{Id: id}).
		Return(nil, errors.New("no user in database"))

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)

	return NewAuthenticate(mockUser, mockAuth, mock.NewSecurityMock(ctrl)), ctrl
}

func empty(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int64)

	if id == id {
		w.WriteHeader(http.StatusOK)
	}
	w.WriteHeader(http.StatusBadRequest)
}

func createExpectedUser() *profile.PublicUser {
	return &profile.PublicUser{Id: id, Email: email, Avatar: avatar}
}
