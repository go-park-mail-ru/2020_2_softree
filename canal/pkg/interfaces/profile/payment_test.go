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
	url := "http://127.0.0.1:8000/api/users"

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
	paymentLogic.EXPECT().ReceiveTransactions(ctx, id).Return(entity.Description{}, createExpectedPayments(), nil)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func TestGetTransactions_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"

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
	paymentLogic.EXPECT().ReceiveTransactions(ctx, id).Return(entity.Description{Status: 500}, nil, errors.New("error"))

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

//
//func TestSetTransaction_FailDecode(t *testing.T) {
//	url := "http://127.0.0.1:8000/api/users"
//
//	req := httptest.NewRequest(http.MethodPost, url, nil)
//	w := httptest.NewRecorder()
//
//	testAuth, ctrl := createSetTransactionFailDecode(t)
//	defer ctrl.Finish()
//
//	createContext(&req)
//	testAuth.SetTransaction(w, req)
//
//	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
//	require.Empty(t, w.Header().Get("Content-Type"))
//	require.Empty(t, w.Body)
//}
//
//func createSetTransactionFailDecode(t *testing.T) (*profileHTTP.Profile, *gomock.Controller) {
//	ctrl := gomock.NewController(t)
//
//	mockUser := profileMock.NewProfileMock(ctrl)
//
//	mockSecurity := mock.NewSecurityMock(ctrl)
//
//	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
//
//	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
//}
//
//func TestSetTransaction_FailCheckFrom400(t *testing.T) {
//	url := "http://127.0.0.1:8000/api/users"
//	body := strings.NewReader(
//		fmt.Sprintf(`{"base": "%s", "currency": "%s", "amount": %f, "sell": "%s"}`, base, curr, amount, sell),
//	)
//
//	req := httptest.NewRequest(http.MethodPost, url, body)
//	w := httptest.NewRecorder()
//
//	testAuth, ctrl := createSetTransactionFailCheckFrom400(t, createContext(&req))
//	defer ctrl.Finish()
//
//	testAuth.SetTransaction(w, req)
//
//	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
//	require.Empty(t, w.Header().Get("Content-Type"))
//	require.Empty(t, w.Body)
//}
//
//func createSetTransactionFailCheckFrom400(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
//	ctrl := gomock.NewController(t)
//
//	mockUser := profileMock.NewProfileMock(ctrl)
//	mockUser.EXPECT().
//		CheckWallet(ctx, createWalletCheck(base)).
//		Return(&profileService.Check{Existence: false}, nil)
//
//	mockSecurity := mock.NewSecurityMock(ctrl)
//
//	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: base}).
//		Return(&currency.Currency{Title: base, Value: value}, nil)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: curr}).
//		Return(&currency.Currency{Title: curr, Value: value}, nil)
//
//	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
//}
//
//func TestSetTransaction_FailCheckFrom500(t *testing.T) {
//	url := "http://127.0.0.1:8000/api/users"
//	body := strings.NewReader(
//		fmt.Sprintf(`{"base": "%s", "currency": "%s", "amount": %f, "sell": "%s"}`, base, curr, amount, sell),
//	)
//
//	req := httptest.NewRequest(http.MethodPost, url, body)
//	w := httptest.NewRecorder()
//
//	testAuth, ctrl := createSetTransactionFailCheckFrom500(t, createContext(&req))
//	defer ctrl.Finish()
//
//	testAuth.SetTransaction(w, req)
//
//	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
//	require.Empty(t, w.Header().Get("Content-Type"))
//	require.Empty(t, w.Body)
//}
//
//func createSetTransactionFailCheckFrom500(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
//	ctrl := gomock.NewController(t)
//
//	mockUser := profileMock.NewProfileMock(ctrl)
//	mockUser.EXPECT().
//		CheckWallet(ctx, createWalletCheck(base)).
//		Return(&profileService.Check{Existence: false}, errors.New("createSetTransactionFailCheckFrom500"))
//
//	mockSecurity := mock.NewSecurityMock(ctrl)
//
//	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: base}).
//		Return(&currency.Currency{Title: base, Value: value}, nil)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: curr}).
//		Return(&currency.Currency{Title: curr, Value: value}, nil)
//
//	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
//}
//
//func TestSetTransaction_FailCurrencyDivFrom(t *testing.T) {
//	url := "http://127.0.0.1:8000/api/users"
//	body := strings.NewReader(
//		fmt.Sprintf(`{"base": "%s", "currency": "%s", "amount": %f, "sell": "%s"}`, base, curr, amount, sell),
//	)
//
//	req := httptest.NewRequest(http.MethodPost, url, body)
//	w := httptest.NewRecorder()
//
//	testAuth, ctrl := createSetTransactionFailCurrencyDivFrom(t, createContext(&req))
//	defer ctrl.Finish()
//
//	testAuth.SetTransaction(w, req)
//
//	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
//	require.Empty(t, w.Header().Get("Content-Type"))
//	require.Empty(t, w.Body)
//}
//
//func createSetTransactionFailCurrencyDivFrom(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
//	ctrl := gomock.NewController(t)
//
//	mockUser := profileMock.NewProfileMock(ctrl)
//
//	mockSecurity := mock.NewSecurityMock(ctrl)
//
//	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: base}).
//		Return(nil, errors.New("createSetTransactionFailCurrencyDivFrom"))
//
//	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
//}
//
//func TestSetTransaction_FailCurrencyDivTo(t *testing.T) {
//	url := "http://127.0.0.1:8000/api/users"
//	body := strings.NewReader(
//		fmt.Sprintf(`{"base": "%s", "currency": "%s", "amount": %f, "sell": "%s"}`, base, curr, amount, sell),
//	)
//
//	req := httptest.NewRequest(http.MethodPost, url, body)
//	w := httptest.NewRecorder()
//
//	testAuth, ctrl := createSetTransactionFailCurrencyDivTo(t, createContext(&req))
//	defer ctrl.Finish()
//
//	testAuth.SetTransaction(w, req)
//
//	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
//	require.Empty(t, w.Header().Get("Content-Type"))
//	require.Empty(t, w.Body)
//}
//
//func createSetTransactionFailCurrencyDivTo(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
//	ctrl := gomock.NewController(t)
//
//	mockUser := profileMock.NewProfileMock(ctrl)
//
//	mockSecurity := mock.NewSecurityMock(ctrl)
//
//	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: base}).
//		Return(&currency.Currency{Title: base, Value: value}, nil)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: curr}).
//		Return(nil, errors.New("createSetTransactionFailCurrencyDivTo"))
//
//	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
//}
//
//func TestSetTransaction_FailGetWallet(t *testing.T) {
//	url := "http://127.0.0.1:8000/api/users"
//	body := strings.NewReader(
//		fmt.Sprintf(`{"base": "%s", "currency": "%s", "amount": %f, "sell": "%s"}`, base, curr, amount, sell),
//	)
//
//	req := httptest.NewRequest(http.MethodPost, url, body)
//	w := httptest.NewRecorder()
//
//	testAuth, ctrl := createSetTransactionFailGetWallet(t, createContext(&req))
//	defer ctrl.Finish()
//
//	testAuth.SetTransaction(w, req)
//
//	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
//	require.Empty(t, w.Header().Get("Content-Type"))
//	require.Empty(t, w.Body)
//}
//
//func createSetTransactionFailGetWallet(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
//	ctrl := gomock.NewController(t)
//
//	mockUser := profileMock.NewProfileMock(ctrl)
//	mockUser.EXPECT().
//		CheckWallet(ctx, createWalletCheck(base)).
//		Return(&profileService.Check{Existence: true}, nil)
//	mockUser.EXPECT().
//		GetWallet(ctx, &profileService.ConcreteWallet{Id: id, Title: base}).
//		Return(nil, errors.New("createSetTransactionFailGetWallet"))
//
//	mockSecurity := mock.NewSecurityMock(ctrl)
//
//	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: base}).
//		Return(&currency.Currency{Title: base, Value: value}, nil)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: curr}).
//		Return(&currency.Currency{Title: curr, Value: value}, nil)
//
//	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
//}
//
//func TestSetTransaction_FailCheckWalletTo(t *testing.T) {
//	url := "http://127.0.0.1:8000/api/users"
//	body := strings.NewReader(
//		fmt.Sprintf(`{"base": "%s", "currency": "%s", "amount": %f, "sell": "%s"}`, base, curr, amount, sell),
//	)
//
//	req := httptest.NewRequest(http.MethodPost, url, body)
//	w := httptest.NewRecorder()
//
//	testAuth, ctrl := createSetTransactionFailCheckWalletTo(t, createContext(&req))
//	defer ctrl.Finish()
//
//	testAuth.SetTransaction(w, req)
//
//	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
//	require.Empty(t, w.Header().Get("Content-Type"))
//	require.Empty(t, w.Body)
//}
//
//func createSetTransactionFailCheckWalletTo(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
//	ctrl := gomock.NewController(t)
//
//	mockUser := profileMock.NewProfileMock(ctrl)
//	mockUser.EXPECT().
//		CheckWallet(ctx, createWalletCheck(base)).
//		Return(&profileService.Check{Existence: true}, nil)
//	mockUser.EXPECT().
//		GetWallet(ctx, &profileService.ConcreteWallet{Id: id, Title: base}).
//		Return(&profileService.Wallet{Title: base, Value: walletSize}, nil)
//	mockUser.EXPECT().
//		CheckWallet(ctx, createWalletCheck(curr)).
//		Return(nil, errors.New("createSetTransactionFailCheckWalletTo"))
//
//	mockSecurity := mock.NewSecurityMock(ctrl)
//
//	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: base}).
//		Return(&currency.Currency{Title: base, Value: value}, nil)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: curr}).
//		Return(&currency.Currency{Title: curr, Value: value}, nil)
//
//	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
//}
//
//func TestSetTransaction_FailCheckWalletToCreateWallet(t *testing.T) {
//	url := "http://127.0.0.1:8000/api/users"
//	body := strings.NewReader(
//		fmt.Sprintf(`{"base": "%s", "currency": "%s", "amount": %f, "sell": "%s"}`, base, curr, amount, sell),
//	)
//
//	req := httptest.NewRequest(http.MethodPost, url, body)
//	w := httptest.NewRecorder()
//
//	testAuth, ctrl := createSetTransactionFailCheckWalletToCreateWallet(t, createContext(&req))
//	defer ctrl.Finish()
//
//	testAuth.SetTransaction(w, req)
//
//	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
//	require.Empty(t, w.Header().Get("Content-Type"))
//	require.Empty(t, w.Body)
//}
//
//func createSetTransactionFailCheckWalletToCreateWallet(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
//	ctrl := gomock.NewController(t)
//
//	mockUser := profileMock.NewProfileMock(ctrl)
//	mockUser.EXPECT().
//		CheckWallet(ctx, createWalletCheck(base)).
//		Return(&profileService.Check{Existence: true}, nil)
//	mockUser.EXPECT().
//		GetWallet(ctx, &profileService.ConcreteWallet{Id: id, Title: base}).
//		Return(&profileService.Wallet{Title: base, Value: walletSize}, nil)
//	mockUser.EXPECT().
//		CheckWallet(ctx, createWalletCheck(curr)).
//		Return(&profileService.Check{Existence: false}, nil)
//	mockUser.EXPECT().
//		CreateWallet(ctx, createWalletCheck(curr)).
//		Return(nil, errors.New("createSetTransactionFailCheckWalletToCreateWallet"))
//
//	mockSecurity := mock.NewSecurityMock(ctrl)
//
//	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: base}).
//		Return(&currency.Currency{Title: base, Value: value}, nil)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: curr}).
//		Return(&currency.Currency{Title: curr, Value: value}, nil)
//
//	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
//}
//
//func TestSetTransaction_FailUpdateFrom(t *testing.T) {
//	url := "http://127.0.0.1:8000/api/users"
//	body := strings.NewReader(
//		fmt.Sprintf(`{"base": "%s", "currency": "%s", "amount": %f, "sell": "%s"}`, base, curr, amount, sell),
//	)
//
//	req := httptest.NewRequest(http.MethodPost, url, body)
//	w := httptest.NewRecorder()
//
//	testAuth, ctrl := createSetTransactionFailUpdateFrom(t, createContext(&req))
//	defer ctrl.Finish()
//
//	testAuth.SetTransaction(w, req)
//
//	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
//	require.Empty(t, w.Header().Get("Content-Type"))
//	require.Empty(t, w.Body)
//}
//
//func createSetTransactionFailUpdateFrom(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
//	ctrl := gomock.NewController(t)
//
//	mockUser := profileMock.NewProfileMock(ctrl)
//	mockUser.EXPECT().
//		CheckWallet(ctx, createWalletCheck(base)).
//		Return(&profileService.Check{Existence: true}, nil)
//	mockUser.EXPECT().
//		GetWallet(ctx, &profileService.ConcreteWallet{Id: id, Title: base}).
//		Return(&profileService.Wallet{Title: base, Value: walletSize}, nil)
//	mockUser.EXPECT().
//		CheckWallet(ctx, createWalletCheck(curr)).
//		Return(&profileService.Check{Existence: true}, nil)
//	mockUser.EXPECT().
//		UpdateWallet(ctx, createWalletToSet(base, -amount)).
//		Return(nil, errors.New("createSetTransactionFailUpdateFrom"))
//
//	mockSecurity := mock.NewSecurityMock(ctrl)
//
//	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: base}).
//		Return(&currency.Currency{Title: base, Value: value}, nil)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: curr}).
//		Return(&currency.Currency{Title: curr, Value: value}, nil)
//
//	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
//}
//
//func TestSetTransaction_FailUpdateTo(t *testing.T) {
//	url := "http://127.0.0.1:8000/api/users"
//	body := strings.NewReader(
//		fmt.Sprintf(`{"base": "%s", "currency": "%s", "amount": %f, "sell": "%s"}`, base, curr, amount, sell),
//	)
//
//	req := httptest.NewRequest(http.MethodPost, url, body)
//	w := httptest.NewRecorder()
//
//	testAuth, ctrl := createSetTransactionFailUpdateTo(t, createContext(&req))
//	defer ctrl.Finish()
//
//	testAuth.SetTransaction(w, req)
//
//	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
//	require.Empty(t, w.Header().Get("Content-Type"))
//	require.Empty(t, w.Body)
//}
//
//func createSetTransactionFailUpdateTo(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
//	ctrl := gomock.NewController(t)
//
//	mockUser := profileMock.NewProfileMock(ctrl)
//	mockUser.EXPECT().
//		CheckWallet(ctx, createWalletCheck(base)).
//		Return(&profileService.Check{Existence: true}, nil)
//	mockUser.EXPECT().
//		GetWallet(ctx, &profileService.ConcreteWallet{Id: id, Title: base}).
//		Return(&profileService.Wallet{Title: base, Value: walletSize}, nil)
//	mockUser.EXPECT().
//		CheckWallet(ctx, createWalletCheck(curr)).
//		Return(&profileService.Check{Existence: true}, nil)
//	mockUser.EXPECT().
//		UpdateWallet(ctx, createWalletToSet(base, -amount)).
//		Return(nil, nil)
//	mockUser.EXPECT().
//		UpdateWallet(ctx, createWalletToSet(curr, amount)).
//		Return(nil, errors.New("createSetTransactionFailUpdateTo"))
//
//	mockSecurity := mock.NewSecurityMock(ctrl)
//
//	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: base}).
//		Return(&currency.Currency{Title: base, Value: value}, nil)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: curr}).
//		Return(&currency.Currency{Title: curr, Value: value}, nil)
//
//	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
//}
//
//func TestSetTransaction_FailAddToPaymentHistory(t *testing.T) {
//	url := "http://127.0.0.1:8000/api/users"
//	body := strings.NewReader(
//		fmt.Sprintf(`{"base": "%s", "currency": "%s", "amount": %f, "sell": "%s"}`, base, curr, amount, sell),
//	)
//
//	req := httptest.NewRequest(http.MethodPost, url, body)
//	w := httptest.NewRecorder()
//
//	testAuth, ctrl := createSetTransactionFailAddToPaymentHistory(t, createContext(&req))
//	defer ctrl.Finish()
//
//	testAuth.SetTransaction(w, req)
//
//	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
//	require.Empty(t, w.Header().Get("Content-Type"))
//	require.Empty(t, w.Body)
//}
//
//func createSetTransactionFailAddToPaymentHistory(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
//	ctrl := gomock.NewController(t)
//
//	mockUser := profileMock.NewProfileMock(ctrl)
//	mockUser.EXPECT().
//		CheckWallet(ctx, createWalletCheck(base)).
//		Return(&profileService.Check{Existence: true}, nil)
//	mockUser.EXPECT().
//		GetWallet(ctx, &profileService.ConcreteWallet{Id: id, Title: base}).
//		Return(&profileService.Wallet{Title: base, Value: walletSize}, nil)
//	mockUser.EXPECT().
//		CheckWallet(ctx, createWalletCheck(curr)).
//		Return(&profileService.Check{Existence: true}, nil)
//	mockUser.EXPECT().
//		UpdateWallet(ctx, createWalletToSet(base, -amount)).
//		Return(nil, nil)
//	mockUser.EXPECT().
//		UpdateWallet(ctx, createWalletToSet(curr, amount)).
//		Return(nil, nil)
//	mockUser.EXPECT().
//		AddToPaymentHistory(ctx, &profileService.AddToHistory{
//			Id: id,
//			Transaction: &profileService.PaymentHistory{Base: base, Currency: curr, Amount: amount, Value: 1, Sell: sell,
//			}}).
//		Return(nil, errors.New("createSetTransactionFailAddToPaymentHistory"))
//
//	mockSecurity := mock.NewSecurityMock(ctrl)
//
//	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: base}).
//		Return(&currency.Currency{Title: base, Value: value}, nil)
//	mockRates.EXPECT().
//		GetLastRate(ctx, &currency.CurrencyTitle{Title: curr}).
//		Return(&currency.Currency{Title: curr, Value: value}, nil)
//
//	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
//}

func createExpectedPayments() entity.Payments {
	return entity.Payments{
		entity.Payment{
			Base: base, Currency: curr, Amount: decimal.NewFromFloat(amount), Value: decimal.NewFromFloat(value), Sell: sell,
		},
	}
}

func createExpectedPayment() entity.Payment {
	return entity.Payment{Base: base, Currency: curr, Amount: decimal.NewFromFloat(amount), Sell: sell, UserId: id}
}

func createWalletCheck(dest string) *profileService.ConcreteWallet {
	return &profileService.ConcreteWallet{Id: id, Title: dest}
}

func createWalletToSet(dest string, val float64) *profileService.ToSetWallet {
	return &profileService.ToSetWallet{Id: id, NewWallet: &profileService.Wallet{Title: dest, Value: val}}
}
