package authorization

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	authMock "server/src/authorization/pkg/infrastructure/mock"
	session "server/src/authorization/pkg/session/gen"
	profileMock "server/src/profile/pkg/infrastructure/mock"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestLogout_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/logout"
	body := strings.NewReader(fmt.Sprintf("{\"email\": %s, \"password\": %s}", email, password))

	req := httptest.NewRequest(http.MethodDelete, url, body)
	w := httptest.NewRecorder()

	ctx := req.Context()
	testAuth, ctrl := createLogoutSuccess(t, &ctx)
	defer ctrl.Finish()

	cookie, _ := CreateCookie()
	cookie.Value = value
	req.AddCookie(&cookie)

	testAuth.Logout(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.Equal(t, w.Result().Cookies()[0].Value, "")
	require.Equal(
		t,
		w.Result().Cookies()[0].Expires,
		time.Date(1973, 1, 1, 0, 0, 0, 0, time.UTC),
	)
}

func TestLogout_FailNoCookie(t *testing.T) {
	url := "http://127.0.0.1:8000/logout"
	body := strings.NewReader(fmt.Sprintf("{\"email\": %s, \"password\": %s}", email, password))

	req := httptest.NewRequest(http.MethodDelete, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createLogoutFailNoCookie(t)
	defer ctrl.Finish()

	testAuth.Logout(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}

func TestLogout_FailDeleteAuth(t *testing.T) {
	url := "http://127.0.0.1:8000/logout"
	body := strings.NewReader(fmt.Sprintf("{\"email\": %s, \"password\": %s}", email, password))

	req := httptest.NewRequest(http.MethodDelete, url, body)
	w := httptest.NewRecorder()

	ctx := req.Context()
	testAuth, ctrl := createLogoutFailDeleteAuth(t, &ctx)
	defer ctrl.Finish()

	cookie, _ := CreateCookie()
	cookie.Value = value
	req.AddCookie(&cookie)

	testAuth.Logout(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func createLogoutSuccess(t *testing.T, ctx *context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().
		Delete(ctx, &session.SessionID{SessionId: value}).
		Return(nil)


	return NewAuthenticate(mockUser, mockAuth), ctrl
}

func createLogoutFailNoCookie(t *testing.T) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)
	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)

	return NewAuthenticate(mockUser, mockAuth), ctrl
}

func createLogoutFailDeleteAuth(t *testing.T, ctx *context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().
		Delete(ctx, &session.SessionID{SessionId: value}).
		Return(errors.New("createLogoutFailDeleteAuth"))

	return NewAuthenticate(mockUser, mockAuth), ctrl
}
