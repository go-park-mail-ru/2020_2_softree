package authorization

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

func TestAuth_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/auth"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createAuthSuccess(t)
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

func TestAuth_FailUnauthorized(t *testing.T) {
	url := "http://127.0.0.1:8000/auth"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createAuthFailUnauthorized(t)
	defer ctrl.Finish()

	testAuth.Auth(w, req)

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

	testAuth.Auth(w, req)

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

	testAuth.Auth(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func createAuthSuccess(t *testing.T) (*Authenticate, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	expectedUser := createExpectedUser("yandex@mail.ru", "str")

	var id uint64 = 1
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().GetUserById(id).Return(expectedUser, nil)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().CheckAuth("value").Return(id, nil)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesCookie := auth.NewToken()
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, *servicesAuth, servicesCookie, servicesLog), ctrl
}

func createAuthFailUnauthorized(t *testing.T) (*Authenticate, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesCookie := auth.NewToken()
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, *servicesAuth, servicesCookie, servicesLog), ctrl
}

func createAuthFailSession(t *testing.T) (*Authenticate, *gomock.Controller) {
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

func createAuthFailUser(t *testing.T) (*Authenticate, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().GetUserById(uint64(1)).Return(entity.User{}, errors.New("no user in database"))

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().CheckAuth("value").Return(uint64(1), nil)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesCookie := auth.NewToken()
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, *servicesAuth, servicesCookie, servicesLog), ctrl
}

func createExpectedUser(email, pass string) (expected entity.User) {
	toSave := entity.User{
		Email:    email,
		Password: pass,
	}
	password, _ := security.MakeShieldedPassword(toSave.Password)
	expected = entity.User{
		ID:       1,
		Email:    toSave.Email,
		Password: password,
	}

	return
}