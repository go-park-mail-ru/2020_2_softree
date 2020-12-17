package profile_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/infrastructure/mock"
	profileHTTP "server/canal/pkg/interfaces/profile"
	"strings"
	"testing"
)

func TestGetWallets_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	testAuth, ctrl := createGetWalletsSuccess(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.GetWallets(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createGetWalletsSuccess(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)
	paymentLogic.EXPECT().ReceiveWallets(ctx, id).Return(entity.Description{}, createExpectedWallets(), nil)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func TestGetWallets_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	testAuth, ctrl := createGetWalletsFail(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.GetWallets(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createGetWalletsFail(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)
	paymentLogic.EXPECT().ReceiveWallets(ctx, id).Return(entity.Description{Status: 500}, entity.Wallets{}, errors.New("error"))

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func TestSetWallets_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(fmt.Sprintf(`{"title": "%s"}`, curr))

	req := httptest.NewRequest(http.MethodGet, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSetWalletsSuccess(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.SetWallet(w, req)

	require.Equal(t, http.StatusCreated, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createSetWalletsSuccess(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)
	paymentLogic.EXPECT().SetWallet(ctx, entity.Wallet{UserId: id, Title: curr}).Return(entity.Description{}, nil)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func TestSetWallets_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(fmt.Sprintf(`{"title": "%s"}`, curr))

	req := httptest.NewRequest(http.MethodGet, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSetWalletsFail(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.SetWallet(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createSetWalletsFail(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)
	paymentLogic.EXPECT().
		SetWallet(ctx, entity.Wallet{UserId: id, Title: curr}).
		Return(entity.Description{Status: 500}, errors.New("error"))

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func TestSetWallets_FailDecode(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSetWalletsFailDecode(t)
	defer ctrl.Finish()

	createContext(&req)
	testAuth.SetWallet(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createSetWalletsFailDecode(t *testing.T) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func createExpectedWallets() entity.Wallets {
	return entity.Wallets{entity.Wallet{Title: curr, Value: decimal.NewFromFloat(value)}}
}
