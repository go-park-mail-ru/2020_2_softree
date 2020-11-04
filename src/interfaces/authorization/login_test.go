package authorization

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"server/src/application"
	"server/src/domain/entity"
	"server/src/infrastructure/log"
	mocks "server/src/infrastructure/mock"
	"strings"
	"testing"
)

func TestLogin_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/login"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createLoginSuccess(t)
	defer ctrl.Finish()

	testAuth.Login(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Result().Cookies())
}

func TestLogin_FailValidation(t *testing.T) {
	url := "http://127.0.0.1:8000/login"
	body := strings.NewReader(`{"email": "yandex.ru", "password": "str"}`)

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
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createLoginFailNoUser(t)
	defer ctrl.Finish()

	testAuth.Login(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	require.Empty(t, w.Body)
	require.Empty(t, w.Header().Get("Content-type"))
}

func TestLogin_FailCreateAuth(t *testing.T) {
	url := "http://127.0.0.1:8000/login"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createLoginFailCreateAuth(t)
	defer ctrl.Finish()

	testAuth.Login(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func createLoginSuccess(t *testing.T) (*Authenticate, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	expectedUser := createExpectedUser("yandex@mail.ru", "str")

	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().GetUserByLogin(expectedUser.Email, "str").Return(expectedUser, nil)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().CreateAuth(expectedUser.ID).Return(http.Cookie{Name: "session_id", Value: "value"}, nil)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func createLoginFailValidation(t *testing.T) (*Authenticate, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func createLoginFailNoUser(t *testing.T) (*Authenticate, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().GetUserByLogin("yandex@mail.ru", "str").Return(entity.User{}, errors.New("no user"))

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func createLoginFailCreateAuth(t *testing.T) (*Authenticate, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	expectedUser := createExpectedUser("yandex@mail.ru", "str")
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().GetUserByLogin(expectedUser.Email, "str").Return(expectedUser, nil)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().CreateAuth(uint64(1)).Return(http.Cookie{}, errors.New("fail to create cookie"))

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, *servicesAuth, servicesLog), ctrl
}
