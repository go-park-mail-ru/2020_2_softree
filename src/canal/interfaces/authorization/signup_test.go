package authorization

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"server/src/canal/application"
	"server/src/canal/domain/entity"
	mocks "server/src/canal/infrastructure/mock"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
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

	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
	require.Equal(t, http.StatusCreated, w.Result().StatusCode)
	require.NotEmpty(t, w.Result().Cookies())
}

func TestSignup_FailUserExist(t *testing.T) {
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(`{"email": "hound@psina.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSignupFailUserExist(t)
	defer ctrl.Finish()

	testAuth.Signup(w, req)

	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestSignup_FailUserExistFailDB(t *testing.T) {
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(`{"email": "hound@psina.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSignupFailUserExistFailDB(t)
	defer ctrl.Finish()

	testAuth.Signup(w, req)

	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
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

func createSignupSuccess(t *testing.T, userToSave entity.User) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	expectedUser := createExpectedUser("hound@psina.ru", "str")
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().CheckExistence(userToSave.Email).Return(false, nil)
	mockUser.EXPECT().SaveUser(userToSave).Times(1).Return(expectedUser, nil)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().CreateAuth(expectedUser.ID).Return(http.Cookie{Name: "session_id"}, nil)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)

	return NewAuthenticate(*servicesDB, *servicesAuth), ctrl
}

func createSignupFail(t *testing.T) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)

	return NewAuthenticate(*servicesDB, *servicesAuth), ctrl
}

func createSignupFailBcrypt(t *testing.T, u entity.User) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().CheckExistence("hound@psina.ru").Return(false, nil)
	mockUser.EXPECT().SaveUser(u).Return(entity.User{}, errors.New("bcrypt: create password hash"))

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)

	return NewAuthenticate(*servicesDB, *servicesAuth), ctrl
}

func createSignupFailUserExist(t *testing.T) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().CheckExistence("hound@psina.ru").Return(true, nil)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)

	return NewAuthenticate(*servicesDB, *servicesAuth), ctrl
}

func createSignupFailUserExistFailDB(t *testing.T) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().CheckExistence("hound@psina.ru").Return(false, errors.New("error"))

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)

	return NewAuthenticate(*servicesDB, *servicesAuth), ctrl
}
