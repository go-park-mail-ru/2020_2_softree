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
	profileService "server/profile/pkg/profile/gen"
	"strings"
	"testing"
)

func TestGetTransactions_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users?period=week"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	testAuth, ctrl := createGetTransactionsSuccess(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.GetTransactions(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createGetTransactionsSuccess(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)
	paymentLogic.EXPECT().ReceiveTransactions(ctx, entity.Income{Id: id, Period: period}).Return(entity.Description{}, createExpectedPayments(), nil)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func TestGetTransactions_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users?period=week"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	testAuth, ctrl := createGetTransactionsFail(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.GetTransactions(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createGetTransactionsFail(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)
	paymentLogic.EXPECT().ReceiveTransactions(ctx, entity.Income{Id: id, Period: period}).Return(entity.Description{Status: 500}, nil, errors.New("error"))

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func TestSetTransaction_FailDecode(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"

	body := strings.NewReader(
		fmt.Sprintf(`{"base": "%s", "currency": "%s", "amount": %f, "sell": "true"}`, base, curr, amount),
	)

	req := httptest.NewRequest(http.MethodPost, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSetTransactionFailDecode(t, req.Context())
	defer ctrl.Finish()

	createContext(&req)
	testAuth.SetTransaction(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createSetTransactionFailDecode(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func TestSetTransaction_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(
		fmt.Sprintf(`{"base": "%s", "currency": "%s", "amount": %f, "sell": %t}`, base, curr, amount, sell),
	)

	req := httptest.NewRequest(http.MethodPost, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSetTransactionSuccess(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.SetTransaction(w, req)

	require.Equal(t, http.StatusCreated, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createSetTransactionSuccess(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)
	paymentLogic.EXPECT().SetTransaction(ctx, createExpectedPayment()).Return(entity.Description{}, nil)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func TestSetTransaction_Fail500(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(
		fmt.Sprintf(`{"base": "%s", "currency": "%s", "amount": %f, "sell": %t}`, base, curr, amount, sell),
	)

	req := httptest.NewRequest(http.MethodPost, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSetTransactionFail500(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.SetTransaction(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createSetTransactionFail500(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)
	paymentLogic.EXPECT().SetTransaction(ctx, createExpectedPayment()).Return(entity.Description{Status: 500}, errors.New("error"))

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func TestSetTransaction_Fail400(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(
		fmt.Sprintf(`{"base": "%s", "currency": "%s", "amount": %f, "sell": %t}`, base, curr, amount, sell),
	)

	req := httptest.NewRequest(http.MethodPost, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSetTransactionFail400(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.SetTransaction(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createSetTransactionFail400(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)
	paymentLogic.EXPECT().SetTransaction(ctx, createExpectedPayment()).Return(entity.Description{Status: 400}, errors.New("error"))

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func createExpectedPayments() entity.Payments {
	return entity.Payments{
		entity.Payment{
			Base: base, Currency: curr, Amount: decimal.NewFromFloat(amount), Value: decimal.NewFromFloat(value), Sell: sell,
		},
	}
}

func createExpectedPayment() entity.Payment {
	return entity.Payment{Base: base, Currency: curr, Amount: decimal.New(1000000000, -6), Sell: sell, UserId: id}
}

func createWalletCheck(dest string) *profileService.ConcreteWallet {
	return &profileService.ConcreteWallet{Id: id, Title: dest}
}

func createWalletToSet(dest string, val float64) *profileService.ToSetWallet {
	return &profileService.ToSetWallet{Id: id, NewWallet: &profileService.Wallet{Title: dest, Value: val}}
}
