package profile

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"server/src/application"
	"server/src/domain/entity"
	"server/src/infrastructure/log"
	mocks "server/src/infrastructure/mock"
	"server/src/infrastructure/security"
	"strings"
	"testing"
)

func TestProfileAvatar_UpdateUserSuccess(t *testing.T) {
	url := "http://127.0.0.1:8000/user"
	body := strings.NewReader(`{"avatar": "fake_image"}`)

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()
	testAuth, ctrl := createUpdateSuccess(t, entity.User{Avatar: "fake_image"})
	defer ctrl.Finish()

	createContext(&req)
	testAuth.UpdateUser(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func TestProfilePassword_UpdateUserSuccess(t *testing.T) {
	url := "http://127.0.0.1:8000/change-password"
	body := strings.NewReader(`{"password": "fake_password"}`)

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createUpdateSuccess(t, entity.User{Password: "fake_password"})
	defer ctrl.Finish()

	createContext(&req)
	testAuth.UpdateUser(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func TestProfile_UpdateUserFail(t *testing.T) {
	url := "http://127.0.0.1:8000/change-password"
	body := strings.NewReader(`{"password": "fake_password"}`)

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createUpdateFail(t, entity.User{Password: "fake_password"})
	defer ctrl.Finish()

	createContext(&req)
	testAuth.UpdateUser(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func TestProfile_AuthSuccess(t *testing.T) {
	url := "http://127.0.0.1:8000/change-password"
	body := strings.NewReader(`{"password": "fake_password"}`)

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createAuthSuccess(t, entity.User{Password: "fake_password"})
	defer ctrl.Finish()

	cookie := http.Cookie{
		Name:  "session_id",
		Value: "value",
	}
	req.AddCookie(&cookie)

	update := testAuth.Auth(testAuth.UpdateUser)
	update(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func TestProfile_AuthFailUnauthorized(t *testing.T) {
	url := "http://127.0.0.1:8000/change-password"
	body := strings.NewReader(`{"password": "fake_password"}`)

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createAuthFailUnauthorized(t)
	defer ctrl.Finish()

	update := testAuth.Auth(testAuth.UpdateUser)
	update(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}

func TestProfile_AuthFailNoSession(t *testing.T) {
	url := "http://127.0.0.1:8000/change-password"
	body := strings.NewReader(`{"password": "fake_password"}`)

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createAuthFailNoSession(t)
	defer ctrl.Finish()

	cookie := http.Cookie{
		Name:  "session_id",
		Value: "value",
	}
	req.AddCookie(&cookie)

	update := testAuth.Auth(testAuth.UpdateUser)
	update(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func createUpdateSuccess(t *testing.T, toUpdate entity.User) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	expectedUser := createExpectedUser()

	var id uint64 = 1
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().UpdateUser(id, toUpdate).Return(expectedUser, nil)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesLog := log.NewLogrusLogger()

	return NewProfile(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func createUpdateFail(t *testing.T, toUpdate entity.User) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	var id uint64 = 1
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().UpdateUser(id, toUpdate).Return(entity.User{}, errors.New("fail to update user"))

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesLog := log.NewLogrusLogger()

	return NewProfile(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func createAuthSuccess(t *testing.T, toUpdate entity.User) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	expectedUser := createExpectedUser()

	var id uint64 = 1
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().UpdateUser(id, toUpdate).Return(expectedUser, nil)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().CheckAuth("value").Return(id, nil)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesLog := log.NewLogrusLogger()

	return NewProfile(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func createAuthFailUnauthorized(t *testing.T) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesLog := log.NewLogrusLogger()

	return NewProfile(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func createAuthFailNoSession(t *testing.T) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mocks.NewUserRepositoryForMock(ctrl)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().CheckAuth("value").Return(uint64(0), errors.New("no session"))

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(mockAuth)
	servicesLog := log.NewLogrusLogger()

	return NewProfile(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func createExpectedUser() (expected entity.User) {
	toSave := entity.User{
		Email: "hound@psina.ru",
		Password: "str",
	}
	password, _ := security.MakeShieldedPassword(toSave.Password)
	expected = entity.User{
		ID: 1,
		Email: toSave.Email,
		Password: password,
		Avatar: "fake_image",
	}

	return
}

func createContext(req **http.Request) {
	ctx := context.WithValue((*req).Context(), "id", uint64(1))
	*req = (*req).Clone(ctx)
}
