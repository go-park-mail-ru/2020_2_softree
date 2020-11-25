package profile_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"server/src/canal/pkg/domain/entity"
	"server/src/canal/pkg/infrastructure/mock"
	profileMock "server/src/profile/pkg/infrastructure/mock"
	profile "server/src/profile/pkg/profile/gen"
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
	password    = "str"
	avatar      = "base64"

	from = "RUB"
	to   = "USD"
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

func createUpdateAvatarSuccess(t *testing.T, ctx context.Context) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		UpdateUserAvatar(ctx, &profile.UpdateFields{Id: id, User: &profile.User{Avatar: avatar}}).
		Return(nil, nil)
	mockUser.EXPECT().
		GetUserById(ctx, &profile.UserID{Id: id}).
		Return(createExpectedUser(), nil)

	mockSecurity := mock.NewSecurityMock(ctrl)

	return NewProfile(mockUser, mockSecurity), ctrl
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

func createUpdateAvatarFail(t *testing.T, ctx context.Context) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		UpdateUserAvatar(ctx, &profile.UpdateFields{Id: id, User: &profile.User{Avatar: avatar}}).
		Return(nil, errors.New("createUpdateAvatarFail"))

	mockSecurity := mock.NewSecurityMock(ctrl)

	return NewProfile(mockUser, mockSecurity), ctrl
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

func createUpdatePasswordSuccess(t *testing.T, ctx context.Context) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckPassword(ctx, &profile.User{OldPassword: oldPassword, NewPassword: newPassword}).
		Return(&profile.Check{Existence: true}, nil)
	mockUser.EXPECT().
		UpdateUserPassword(id, &profile.UpdateFields{
			Id: id,
			User: &profile.User{OldPassword: oldPassword, NewPassword: newPassword,
			}}).
		Return(nil)
	mockUser.EXPECT().
		GetUserById(ctx, &profile.UserID{Id: id}).
		Return(createExpectedUser(), nil)

	mockSecurity := mock.NewSecurityMock(ctrl)

	return NewProfile(mockUser, mockSecurity), ctrl
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

func createUpdatePasswordFail(t *testing.T, ctx context.Context) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckPassword(ctx, &profile.User{OldPassword: oldPassword, NewPassword: newPassword}).
		Return(&profile.Check{Existence: false}, errors.New("createUpdatePasswordFail"))

	mockSecurity := mock.NewSecurityMock(ctrl)

	return NewProfile(mockUser, mockSecurity), ctrl
}

func TestGetUser_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/users"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	testAuth, ctrl := createGetUserSuccess(t, createContext(&req))
	defer ctrl.Finish()

	createContext(&req)
	testAuth.GetUser(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createGetUserSuccess(t *testing.T, ctx context.Context) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		GetUserById(ctx, &profile.UserID{Id: id}).
		Return(createExpectedUser(), nil)

	mockSecurity := mock.NewSecurityMock(ctrl)

	return NewProfile(mockUser, mockSecurity), ctrl
}

func TestGetUser_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/rates"
	body := strings.NewReader(`{"avatar": "QmFzZTY0"}`)

	req := httptest.NewRequest(http.MethodGet, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createGetUserFail(t, createContext(&req))
	defer ctrl.Finish()

	createContext(&req)
	testAuth.GetUser(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createGetUserFail(t *testing.T, ctx context.Context) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		GetUserById(ctx, &profile.UserID{Id: id}).
		Return(nil, errors.New("createGetUserFail"))

	mockSecurity := mock.NewSecurityMock(ctrl)

	return NewProfile(mockUser, mockSecurity), ctrl
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

func createGetUserWatchlistSuccess(t *testing.T, ctx context.Context) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		GetUserWatchlist(ctx, &profile.UserID{Id: id}).
		Return(createExpectedCurrencies(), nil)

	mockSecurity := mock.NewSecurityMock(ctrl)

	return NewProfile(mockUser, mockSecurity), ctrl
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

func createGetUserWatchlistFail(t *testing.T, ctx context.Context) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		GetUserWatchlist(ctx, &profile.UserID{Id: id}).
		Return(nil, errors.New("createGetUserWatchlistFail"))

	mockSecurity := mock.NewSecurityMock(ctrl)

	return NewProfile(mockUser, mockSecurity), ctrl
}

func createExpectedUser() *profile.PublicUser {
	return &profile.PublicUser{Id: id, Email: email, Avatar: avatar}
}

func createContext(req **http.Request) context.Context {
	ctx := context.WithValue((*req).Context(), "id", int64(id))
	*req = (*req).Clone(ctx)
	return ctx
}

func createExpectedCurrencies() *profile.Currencies {
	return &profile.Currencies{Currencies: []*profile.Currency{{Base: to, Title: from}}}
}
