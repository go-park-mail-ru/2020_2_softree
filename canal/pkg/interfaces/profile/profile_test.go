package profile_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/infrastructure/mock"
	profileHTTP "server/canal/pkg/interfaces/profile"
	currencyMock "server/currency/pkg/infrastructure/mock"
	profileMock "server/profile/pkg/infrastructure/mock"
	profileService "server/profile/pkg/profile/gen"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	id          = 1
	email       = "hound@psina.ru"
	oldPassword = "old"
	newPassword = "new"
	avatar      = "base64"

	walletSize = 100000.0
	value      = 79.7
	amount     = 1000.0
	from       = "RUB"
	to         = "USD"
)

func TestUpdateUserAvatar_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(fmt.Sprintf(`{"avatar": "%s"}`, avatar))

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createUpdateAvatarSuccess(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.UpdateUserAvatar(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createUpdateAvatarSuccess(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		UpdateUserAvatar(ctx, &profileService.UpdateFields{Id: id, User: &profileService.User{Id: id, Avatar: avatar}}).
		Return(&profileService.Empty{}, nil)
	mockUser.EXPECT().
		GetUserById(ctx, &profileService.UserID{Id: id}).
		Return(createExpectedUser(), nil)

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestUpdateUserAvatar_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(fmt.Sprintf(`{"avatar": "%s"}`, avatar))

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createUpdateAvatarFail(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.UpdateUserAvatar(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createUpdateAvatarFail(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		UpdateUserAvatar(ctx, &profileService.UpdateFields{Id: id, User: &profileService.User{Id: id, Avatar: avatar}}).
		Return(&profileService.Empty{}, errors.New("createUpdateAvatarFail"))

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestUpdateUserPassword_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users/change-password"
	body := strings.NewReader(fmt.Sprintf(`{"old_password": "%s", "new_password": "%s"}`, oldPassword, newPassword))

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createUpdatePasswordSuccess(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.UpdateUserPassword(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createUpdatePasswordSuccess(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckPassword(ctx, &profileService.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword}).
		Return(&profileService.Check{Existence: true}, nil)
	mockUser.EXPECT().
		UpdateUserPassword(
			ctx,
			&profileService.UpdateFields{Id: id, User: &profileService.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword}},
		).
		Return(&profileService.Empty{}, nil)
	mockUser.EXPECT().
		GetUserById(ctx, &profileService.UserID{Id: id}).
		Return(createExpectedUser(), nil)

	mockSecurity := mock.NewSecurityMock(ctrl)
	mockSecurity.EXPECT().MakeShieldedPassword(newPassword).Return(newPassword, nil)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestUpdateUserPassword_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users/change-password"
	body := strings.NewReader(fmt.Sprintf(`{"old_password": "%s", "new_password": "%s"}`, oldPassword, newPassword))

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createUpdatePasswordFail(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.UpdateUserPassword(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func createUpdatePasswordFail(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckPassword(ctx, &profileService.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword}).
		Return(&profileService.Check{Existence: false}, errors.New("createUpdatePasswordFail"))

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestGetUser_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/users"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	testAuth, ctrl := createGetUserSuccess(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.GetUser(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createGetUserSuccess(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		GetUserById(ctx, &profileService.UserID{Id: id}).
		Return(createExpectedUser(), nil)

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestGetUser_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/rates"
	body := strings.NewReader(`{"avatar": "QmFzZTY0"}`)

	req := httptest.NewRequest(http.MethodGet, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createGetUserFail(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.GetUser(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createGetUserFail(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		GetUserById(ctx, &profileService.UserID{Id: id}).
		Return(&profileService.PublicUser{}, errors.New("createGetUserFail"))

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestGetUserWatchlist_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/rates"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	testAuth, ctrl := createGetUserWatchlistSuccess(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.GetUserWatchlist(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createGetUserWatchlistSuccess(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		GetUserWatchlist(ctx, &profileService.UserID{Id: id}).
		Return(createExpectedCurrencies(), nil)

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestGetUserWatchlist_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/rates"
	body := strings.NewReader(`{"avatar": "QmFzZTY0"}`)

	req := httptest.NewRequest(http.MethodGet, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createGetUserWatchlistFail(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.GetUserWatchlist(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createGetUserWatchlistFail(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		GetUserWatchlist(ctx, &profileService.UserID{Id: id}).
		Return(nil, errors.New("createGetUserWatchlistFail"))

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func createExpectedUser() *profileService.PublicUser {
	return &profileService.PublicUser{Id: id, Email: email, Avatar: avatar}
}

func createContext(req **http.Request) context.Context {
	ctx := context.WithValue((*req).Context(), entity.UserIdKey, int64(id))
	*req = (*req).Clone(ctx)
	return ctx
}

func createExpectedCurrencies() *profileService.Currencies {
	return &profileService.Currencies{Currencies: []*profileService.Currency{{Base: to, Title: from}}}
}
