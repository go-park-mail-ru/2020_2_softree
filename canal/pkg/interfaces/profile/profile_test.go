package profile_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/infrastructure/mock"
	profileHTTP "server/canal/pkg/interfaces/profile"
	"strings"
	"testing"
)

const (
	id          = int64(1)
	email       = "hound@psina.ru"
	oldPassword = "old"
	newPassword = "new"
	avatar      = "base64"

	period = "week"
	value  = 79.7
	amount = 1000.0
	base   = "RUB"
	curr   = "USD"
	sell   = false
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

	profileLogic := mock.NewMockProfileLogic(ctrl)
	profileLogic.EXPECT().
		UpdateAvatar(ctx, entity.User{Id: id, Avatar: avatar}).
		Return(entity.Description{}, createExpectedPublicUser(), nil)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
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

	profileLogic := mock.NewMockProfileLogic(ctrl)
	profileLogic.EXPECT().
		UpdateAvatar(ctx, entity.User{Id: id, Avatar: avatar}).
		Return(entity.Description{Status: 500}, entity.PublicUser{}, errors.New("error"))

	paymentLogic := mock.NewMockPaymentLogic(ctrl)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func TestUpdateUserAvatar_FailGetBody(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(fmt.Sprintf(`{"avatar": %d}`, id))

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createUpdateAvatarFailGetBody(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.UpdateUserAvatar(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createUpdateAvatarFailGetBody(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
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

	profileLogic := mock.NewMockProfileLogic(ctrl)
	profileLogic.EXPECT().
		UpdatePassword(ctx, entity.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword}).
		Return(entity.Description{}, createExpectedPublicUser(), nil)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
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

	profileLogic := mock.NewMockProfileLogic(ctrl)
	profileLogic.EXPECT().
		UpdatePassword(ctx, entity.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword}).
		Return(entity.Description{Status: 500}, entity.PublicUser{}, errors.New("error"))

	paymentLogic := mock.NewMockPaymentLogic(ctrl)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func TestUpdateUserPassword_FailGetBody(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users/change-password"
	body := strings.NewReader(fmt.Sprintf(`{"old_password": %d, "new_password": "%s"}`, id, newPassword))

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createUpdatePasswordFailGetBody(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.UpdateUserPassword(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func createUpdatePasswordFailGetBody(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func TestUpdateUserPassword_FailErrorJson(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users/change-password"
	body := strings.NewReader(fmt.Sprintf(`{"old_password": "%s", "new_password": "%s"}`, oldPassword, newPassword))

	req := httptest.NewRequest(http.MethodPut, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createUpdatePasswordFailErrorJson(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.UpdateUserPassword(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createUpdatePasswordFailErrorJson(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)
	profileLogic.EXPECT().
		UpdatePassword(ctx, entity.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword}).
		Return(entity.Description{ErrorJSON: entity.ErrorJSON{NotEmpty: true, NonFieldError: []string{"error"}}}, entity.PublicUser{}, nil)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
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

	profileLogic := mock.NewMockProfileLogic(ctrl)
	profileLogic.EXPECT().ReceiveUser(ctx, id).Return(entity.Description{}, createExpectedPublicUser(), nil)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
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

	profileLogic := mock.NewMockProfileLogic(ctrl)
	profileLogic.EXPECT().ReceiveUser(ctx, id).Return(entity.Description{Status: 500}, entity.PublicUser{}, errors.New("error"))

	paymentLogic := mock.NewMockPaymentLogic(ctrl)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
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

	profileLogic := mock.NewMockProfileLogic(ctrl)
	profileLogic.EXPECT().ReceiveWatchlist(ctx, id).Return(entity.Description{}, createExpectedCurrencies(), nil)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func TestGetUserWatchlist_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/rates"
	body := strings.NewReader(`{"avatar": "QmFzZTY0"}`)

	req := httptest.NewRequest(http.MethodGet, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createGetUserWatchlistFail(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.GetUserWatchlist(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createGetUserWatchlistFail(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)
	profileLogic.EXPECT().ReceiveWatchlist(ctx, id).Return(entity.Description{Status: 500}, entity.Currencies{}, errors.New("error"))

	paymentLogic := mock.NewMockPaymentLogic(ctrl)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func createExpectedPublicUser() entity.PublicUser {
	return entity.PublicUser{Id: id, Email: email, Avatar: avatar}
}

func createContext(req **http.Request) context.Context {
	ctx := context.WithValue((*req).Context(), entity.UserIdKey, id)
	*req = (*req).Clone(ctx)
	return ctx
}

func createExpectedCurrencies() entity.Currencies {
	return entity.Currencies{entity.Currency{Base: curr, Title: base}}
}
