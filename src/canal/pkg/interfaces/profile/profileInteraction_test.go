package profile

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"server/src/canal/pkg/application"
	"server/src/canal/pkg/domain/entity"
	mocks "server/src/canal/pkg/infrastructure/mock"
	"server/src/canal/pkg/infrastructure/security"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserAvatar_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(`{"avatar": "QmFzZTY0"}`)

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createUpdateAvatarSuccess(t, entity.User{Avatar: "QmFzZTY0"})
	defer ctrl.Finish()

	createContext(&req)
	testAuth.UpdateUserAvatar(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createUpdateAvatarSuccess(t *testing.T, toUpdate entity.User) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	var id uint64 = 1
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().UpdateUserAvatar(id, toUpdate).Return(nil)
	mockUser.EXPECT().GetUserById(id).Return(createExpectedUser(), nil)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	mockHistory := mocks.NewPaymentHistoryRepositoryForMock(ctrl)
	mockWallet := mocks.NewWalletRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser, mockHistory, mockWallet)
	servicesAuth := application.NewUserAuth(mockAuth)

	return NewProfile(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func TestUpdateUserAvatar_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(`{"avatar": "QmFzZTY0"}`)

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createUpdateAvatarFail(t, entity.User{Avatar: "QmFzZTY0"})
	defer ctrl.Finish()

	createContext(&req)
	testAuth.UpdateUserAvatar(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createUpdateAvatarFail(t *testing.T, toUpdate entity.User) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().UpdateUserAvatar(uint64(1), toUpdate).Return(errors.New("error"))

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	mockHistory := mocks.NewPaymentHistoryRepositoryForMock(ctrl)
	mockWallet := mocks.NewWalletRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser, mockHistory, mockWallet)
	servicesAuth := application.NewUserAuth(mockAuth)

	return NewProfile(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func TestUpdateUserPassword_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users/change-password"
	body := strings.NewReader(`{"old_password": "fake_password", "new_password": "str"}`)

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createUpdatePasswordSuccess(t, entity.User{OldPassword: "fake_password", NewPassword: "str"})
	defer ctrl.Finish()

	createContext(&req)
	testAuth.UpdateUserPassword(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createUpdatePasswordSuccess(t *testing.T, toUpdate entity.User) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	var id uint64 = 1
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().CheckPassword(id, toUpdate.OldPassword).Return(true, nil)
	mockUser.EXPECT().UpdateUserPassword(id, toUpdate).Return(nil)
	mockUser.EXPECT().GetUserById(id).Return(createExpectedUser(), nil)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	mockHistory := mocks.NewPaymentHistoryRepositoryForMock(ctrl)
	mockWallet := mocks.NewWalletRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser, mockHistory, mockWallet)
	servicesAuth := application.NewUserAuth(mockAuth)

	return NewProfile(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func TestOldPassError(t *testing.T) {
	w := httptest.NewRecorder()

	testAuth, ctrl := createAuthFailUnauthorized(t)
	defer ctrl.Finish()

	testAuth.createOldPassError(w)
	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createAuthFailUnauthorized(t *testing.T) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	mockHistory := mocks.NewPaymentHistoryRepositoryForMock(ctrl)
	mockWallet := mocks.NewWalletRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser, mockHistory, mockWallet)
	servicesAuth := application.NewUserAuth(mockAuth)

	return NewProfile(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func TestUpdateUserPassword_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users/change-password"
	body := strings.NewReader(`{"old_password": "fake_password", "new_password": "str"}`)

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createUpdatePasswordFail(t, entity.User{OldPassword: "fake_password", NewPassword: "str"})
	defer ctrl.Finish()

	createContext(&req)
	testAuth.UpdateUserPassword(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func createUpdatePasswordFail(t *testing.T, toUpdate entity.User) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	var id uint64 = 1
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().CheckPassword(id, toUpdate.OldPassword).Return(false, nil)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	mockHistory := mocks.NewPaymentHistoryRepositoryForMock(ctrl)
	mockWallet := mocks.NewWalletRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser, mockHistory, mockWallet)
	servicesAuth := application.NewUserAuth(mockAuth)

	return NewProfile(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func TestAuth_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/change-password"
	body := strings.NewReader(`{"avatar": "QmFzZTY0"}`)

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createAuthSuccess(t, entity.User{Avatar: "QmFzZTY0"})
	defer ctrl.Finish()

	cookie := http.Cookie{
		Name:  "session_id",
		Value: "value",
	}
	req.AddCookie(&cookie)
	createContext(&req)

	update := testAuth.Auth(testAuth.UpdateUserAvatar)
	update(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createAuthSuccess(t *testing.T, toUpdate entity.User) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	var id uint64 = 1
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().UpdateUserAvatar(id, toUpdate).Return(nil)
	mockUser.EXPECT().GetUserById(id).Return(createExpectedUser(), nil)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().CheckAuth("value").Return(id, nil)

	mockHistory := mocks.NewPaymentHistoryRepositoryForMock(ctrl)
	mockWallet := mocks.NewWalletRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser, mockHistory, mockWallet)
	servicesAuth := application.NewUserAuth(mockAuth)

	return NewProfile(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func TestGetUser_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/rates"
	body := strings.NewReader(`{"avatar": "QmFzZTY0"}`)

	req := httptest.NewRequest(http.MethodGet, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createGetUserSuccess(t)
	defer ctrl.Finish()

	createContext(&req)
	testAuth.GetUser(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createGetUserSuccess(t *testing.T) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().GetUserById(uint64(1)).Return(createExpectedUser(), nil)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	mockHistory := mocks.NewPaymentHistoryRepositoryForMock(ctrl)
	mockWallet := mocks.NewWalletRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser, mockHistory, mockWallet)
	servicesAuth := application.NewUserAuth(mockAuth)

	return NewProfile(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func TestGetUser_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/rates"
	body := strings.NewReader(`{"avatar": "QmFzZTY0"}`)

	req := httptest.NewRequest(http.MethodGet, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createGetUserFail(t)
	defer ctrl.Finish()

	createContext(&req)
	testAuth.GetUser(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createGetUserFail(t *testing.T) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().GetUserById(uint64(1)).Return(createExpectedUser(), errors.New("error"))

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	mockHistory := mocks.NewPaymentHistoryRepositoryForMock(ctrl)
	mockWallet := mocks.NewWalletRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser, mockHistory, mockWallet)
	servicesAuth := application.NewUserAuth(mockAuth)

	return NewProfile(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func TestGetUserWatchlist_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/rates"
	body := strings.NewReader(`{"avatar": "QmFzZTY0"}`)

	req := httptest.NewRequest(http.MethodGet, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createGetUserWatchlistSuccess(t)
	defer ctrl.Finish()

	createContext(&req)
	testAuth.GetUserWatchlist(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createGetUserWatchlistSuccess(t *testing.T) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().GetUserWatchlist(uint64(1)).Return([]entity.Currency{{Title: "USD"}}, nil)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	mockHistory := mocks.NewPaymentHistoryRepositoryForMock(ctrl)
	mockWallet := mocks.NewWalletRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser, mockHistory, mockWallet)
	servicesAuth := application.NewUserAuth(mockAuth)

	return NewProfile(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func TestGetUserWatchlist_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/rates"
	body := strings.NewReader(`{"avatar": "QmFzZTY0"}`)

	req := httptest.NewRequest(http.MethodGet, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createGetUserWatchlistFail(t)
	defer ctrl.Finish()

	createContext(&req)
	testAuth.GetUserWatchlist(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createGetUserWatchlistFail(t *testing.T) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().GetUserWatchlist(uint64(1)).Return(nil, errors.New("error"))

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)

	mockHistory := mocks.NewPaymentHistoryRepositoryForMock(ctrl)
	mockWallet := mocks.NewWalletRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser, mockHistory, mockWallet)
	servicesAuth := application.NewUserAuth(mockAuth)

	return NewProfile(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func TestAuth_FailUnauthorized(t *testing.T) {
	url := "http://127.0.0.1:8000/change-password"
	body := strings.NewReader(`{"password": "fake_password"}`)

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createAuthFailUnauthorized(t)
	defer ctrl.Finish()

	update := testAuth.Auth(testAuth.UpdateUserAvatar)
	update(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}

func TestAuth_FailNoSession(t *testing.T) {
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

	update := testAuth.Auth(testAuth.UpdateUserAvatar)
	update(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func createAuthFailNoSession(t *testing.T) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mocks.NewUserRepositoryForMock(ctrl)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().CheckAuth("value").Return(uint64(0), errors.New("no session"))

	mockHistory := mocks.NewPaymentHistoryRepositoryForMock(ctrl)
	mockWallet := mocks.NewWalletRepositoryForMock(ctrl)

	servicesDB := application.NewUserApp(mockUser, mockHistory, mockWallet)
	servicesAuth := application.NewUserAuth(mockAuth)

	return NewProfile(*servicesDB, *servicesAuth, servicesLog), ctrl
}

func createExpectedUser() (expected entity.User) {
	toSave := entity.User{
		Email:    "hound@psina.ru",
		Password: "str",
	}
	password, _ := security.MakeShieldedPassword(toSave.Password)
	expected = entity.User{
		ID:       1,
		Email:    toSave.Email,
		Password: password,
		Avatar:   "fake_image",
	}

	return
}

func createContext(req **http.Request) {
	ctx := context.WithValue((*req).Context(), "id", uint64(1))
	*req = (*req).Clone(ctx)
}
