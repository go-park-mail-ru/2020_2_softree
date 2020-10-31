package pureArchAuth

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"server/src/application"
	"server/src/domain/entity"
	"server/src/infrastructure/auth"
	"server/src/infrastructure/log"
	mocks "server/src/infrastructure/mock"
	"server/src/infrastructure/security"
	"strings"
	"testing"
)

func TestAuthSuccess(t *testing.T) {
	url := "http://127.0.0.1:8000/auth"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createTestAuthAuthenticateSuccess(t)
	defer ctrl.Finish()

	cookie := http.Cookie{
		Name:  "session_id",
		Value: "value",
	}
	req.AddCookie(&cookie)

	testAuth.Auth(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-type"))
	require.NotEmpty(t, w.Body)
}

func TestAuthFail_Unauthorized(t *testing.T) {
	url := "http://127.0.0.1:8000/auth"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createTestAuthAuthenticateFailUnauthorized(t)
	defer ctrl.Finish()

	testAuth.Auth(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}

func TestAuthFail_NoSession(t *testing.T) {
	url := "http://127.0.0.1:8000/auth"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createTestAuthAuthenticateFailSession(t)
	defer ctrl.Finish()

	cookie := http.Cookie{
		Name:  "session_id",
		Value: "value",
	}
	req.AddCookie(&cookie)

	testAuth.Auth(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestAuthFail_NoUser(t *testing.T) {
	url := "http://127.0.0.1:8000/auth"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createTestAuthAuthenticateFailUser(t)
	defer ctrl.Finish()

	cookie := http.Cookie{
		Name:  "session_id",
		Value: "value",
	}
	req.AddCookie(&cookie)

	testAuth.Auth(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func createTestAuthAuthenticateSuccess(t *testing.T) (*Authenticate, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	expectedUser := createExpectedUser()

	var id uint64 = 1
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().GetUser(id).Return(expectedUser, nil)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().CheckAuth("value").Return(id, nil)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesCookie := auth.NewToken()
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, *servicesAuth, servicesCookie, servicesLog), ctrl
}

func createTestAuthAuthenticateFailUnauthorized(t *testing.T) (*Authenticate, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesCookie := auth.NewToken()
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, *servicesAuth, servicesCookie, servicesLog), ctrl
}

func createTestAuthAuthenticateFailSession(t *testing.T) (*Authenticate, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mocks.NewUserRepositoryForMock(ctrl)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().CheckAuth("value").Return(uint64(0), errors.New("no session"))

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesCookie := auth.NewToken()
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, *servicesAuth, servicesCookie, servicesLog), ctrl
}

func createTestAuthAuthenticateFailUser(t *testing.T) (*Authenticate, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().GetUser(uint64(1)).Return(entity.User{}, errors.New("No user in database"))

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().CheckAuth("value").Return(uint64(1), nil)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesCookie := auth.NewToken()
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, *servicesAuth, servicesCookie, servicesLog), ctrl
}

func createExpectedUser() (expected entity.User) {
	toSave := entity.User{
		Email:    "yandex@mail.ru",
		Password: "str",
	}
	password, _ := security.MakeShieldedPassword(toSave.Password)
	expected = entity.User{
		ID:       1,
		Email:    toSave.Email,
		Password: password,
	}

	return
}
