package authorization

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	authMock "server/authorization/pkg/infrastructure/mock"
	session "server/authorization/pkg/session/gen"
	"server/canal/pkg/infrastructure/mock"
	profileMock "server/profile/pkg/infrastructure/mock"
	profile "server/profile/pkg/profile/gen"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestLogin_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/sessions"
	body := strings.NewReader(fmt.Sprintf("{\"email\": \"%s\", \"password\": \"%s\"}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	ctx := req.Context()
	testAuth, ctrl := createLoginSuccess(t, ctx)
	defer ctrl.Finish()

	testAuth.Login(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Result().Cookies())
}

func TestLogin_FailValidation(t *testing.T) {
	url := "http://127.0.0.1:8000/login"
	body := strings.NewReader(fmt.Sprintf("{\"email\": \"%s\", \"password\": \"%s\"}", "", password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createLoginFailValidation(t)
	defer ctrl.Finish()

	testAuth.Login(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	require.NotEmpty(t, w.Body)
	require.NotEmpty(t, w.Header().Get("Content-type"))
}

func TestLogin_FailNoUser(t *testing.T) {
	url := "http://127.0.0.1:8000/login"
	body := strings.NewReader(fmt.Sprintf("{\"email\": \"%s\", \"password\": \"%s\"}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	ctx := req.Context()
	testAuth, ctrl := createLoginFailNoUser(t, ctx)
	defer ctrl.Finish()

	testAuth.Login(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	require.NotEmpty(t, w.Body)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
}

func TestLogin_FailCreateAuth(t *testing.T) {
	url := "http://127.0.0.1:8000/login"
	body := strings.NewReader(fmt.Sprintf("{\"email\": \"%s\", \"password\": \"%s\"}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	ctx := req.Context()
	testAuth, ctrl := createLoginFailCreateAuth(t, ctx)
	defer ctrl.Finish()

	testAuth.Login(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func createLoginSuccess(t *testing.T, ctx context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email, Password: password}).
		Return(&profile.Check{Existence: true}, nil)
	mockUser.EXPECT().
		GetUserByLogin(ctx, &profile.User{Email: email, Password: password}).
		Return(createExpectedUser(), nil)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().
		Create(ctx, &session.UserID{Id: id}).
		Return(&session.Session{Id: id, SessionId: value}, nil)

	return NewAuthenticate(mockUser, mockAuth, mock.NewSecurityMock(ctrl)), ctrl
}

func createLoginFailValidation(t *testing.T) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)
	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)

	return NewAuthenticate(mockUser, mockAuth, mock.NewSecurityMock(ctrl)), ctrl
}

func createLoginFailNoUser(t *testing.T, ctx context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email, Password: password}).
		Return(&profile.Check{Existence: false}, nil)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)

	return NewAuthenticate(mockUser, mockAuth, mock.NewSecurityMock(ctrl)), ctrl
}

func createLoginFailCreateAuth(t *testing.T, ctx context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email, Password: password}).
		Return(&profile.Check{Existence: true}, nil)
	mockUser.
		EXPECT().GetUserByLogin(ctx, &profile.User{Email: email, Password: password}).
		Return(createExpectedUser(), nil)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().
		Create(ctx, &session.UserID{Id: id}).
		Return(nil, errors.New("fail to create cookie"))

	return NewAuthenticate(mockUser, mockAuth, mock.NewSecurityMock(ctrl)), ctrl
}
