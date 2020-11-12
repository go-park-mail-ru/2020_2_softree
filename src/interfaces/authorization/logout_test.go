package authorization

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"server/src/application"
	"server/src/infrastructure/log"
	mocks "server/src/infrastructure/mock"
	"strings"
	"testing"
	"time"
)

func TestLogout_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/logout"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	testAuth, ctrl := createLogoutSuccess(t)
	defer ctrl.Finish()

	cookie, _ := testAuth.auth.CreateCookie()
	cookie.Value = "value"
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
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	testAuth, ctrl := createLogoutFailNoCookie(t)
	defer ctrl.Finish()

	testAuth.Logout(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}

func TestLogout_FailDeleteAuth(t *testing.T) {
	url := "http://127.0.0.1:8000/logout"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createLogoutFailDeleteAuth(t)
	defer ctrl.Finish()

	cookie, _ := testAuth.auth.CreateCookie()
	cookie.Value = "value"
	req.AddCookie(&cookie)

	testAuth.Logout(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func createLogoutSuccess(t *testing.T) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mocks.NewUserRepositoryForMock(ctrl)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().DeleteAuth("value").Return(nil)
	mockAuth.EXPECT().CreateCookie().Return(http.Cookie{Name: "session_id", Value: "value"}, nil).Times(2)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func createLogoutFailNoCookie(t *testing.T) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func createLogoutFailDeleteAuth(t *testing.T) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mocks.NewUserRepositoryForMock(ctrl)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().DeleteAuth("value").Return(errors.New("delete auth"))
	mockAuth.EXPECT().CreateCookie().Return(http.Cookie{Name: "session_id", Value: "value"}, nil)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, *servicesAuth, servicesLog), ctrl
}
