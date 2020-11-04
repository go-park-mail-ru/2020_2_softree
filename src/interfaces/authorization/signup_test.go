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
	"strings"
	"testing"
)

func TestSignup_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(`{"email": "hound@psina.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSignupSuccess(t, entity.User{
		Email:    "hound@psina.ru",
		Password: "str",
	})
	defer ctrl.Finish()

	testAuth.Signup(w, req)

	require.Empty(t, w.Header().Get("Content-type"))
	require.Empty(t, w.Body)
	require.Equal(t, http.StatusCreated, w.Result().StatusCode)
	require.NotEmpty(t, w.Result().Cookies())
}

func TestSignup_FailEmail(t *testing.T) {
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(`{"email": "hound.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSignupFail(t)
	defer ctrl.Finish()

	testAuth.Signup(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	require.NotEmpty(t, w.Body)
	require.NotEmpty(t, w.Header().Get("Content-type"))
}

func TestSignup_FailEmptyPassword(t *testing.T) {
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(`{"email": "hound@psina.ru"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSignupFail(t)
	defer ctrl.Finish()

	testAuth.Signup(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	require.NotEmpty(t, w.Body)
	require.NotEmpty(t, w.Header().Get("Content-type"))
}

func TestSignup_FailBcrypt(t *testing.T) {
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(`{"email": "hound@psina.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSignupFailBcrypt(t, entity.User{
		Email:    "hound@psina.ru",
		Password: "str",
	})
	defer ctrl.Finish()

	testAuth.Signup(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func createSignupSuccess(t *testing.T, userToSave entity.User) (*Authenticate, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	expectedUser := createExpectedUser("hound@psina.ru", "str")
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().SaveUser(userToSave).Times(1).Return(expectedUser, nil)

	memAuth := auth.NewMemAuth()

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(memAuth)
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func createSignupFail(t *testing.T) (*Authenticate, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func createSignupFailBcrypt(t *testing.T, u entity.User) (*Authenticate, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().SaveUser(u).Return(entity.User{}, errors.New("bcrypt: create password hash"))

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, *servicesAuth, servicesLog), ctrl
}
