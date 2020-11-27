package profile_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"server/src/canal/pkg/infrastructure/mock"
	profileHTTP "server/src/canal/pkg/interfaces/profile"
	currency "server/src/currency/pkg/currency/gen"
	currencyMock "server/src/currency/pkg/infrastructure/mock"
	profileMock "server/src/profile/pkg/infrastructure/mock"
	profileService "server/src/profile/pkg/profile/gen"
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

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		GetAllPaymentHistory(ctx, &profileService.UserID{Id: id}).
		Return(createExpectedHistory(), nil)

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
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

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		GetAllPaymentHistory(ctx, &profileService.UserID{Id: id}).
		Return(nil, errors.New("createGetTransactionsFail"))

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestSetTransaction_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(fmt.Sprintf(`{"from": "%s", "to": "%s", "amount": %f}`, from, to, amount))

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

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckWallet(ctx, createWalletCheck(from)).
		Return(&profileService.Check{Existence: true}, nil)
	mockUser.EXPECT().
		GetWallet(ctx, &profileService.ConcreteWallet{Id: id, Title: to}).
		Return(&profileService.Wallet{Title: to, Value: walletSize}, nil)
	mockUser.EXPECT().
		CheckWallet(ctx, createWalletCheck(to)).
		Return(&profileService.Check{Existence: true}, nil)
	mockUser.EXPECT().
		UpdateWallet(ctx, createWalletToSet(from, -amount)).
		Return(nil, nil)
	mockUser.EXPECT().
		UpdateWallet(ctx, createWalletToSet(to, amount)).
		Return(nil, nil)
	mockUser.EXPECT().
		AddToPaymentHistory(ctx, &profileService.AddToHistory{Id: id, Transaction: &profileService.PaymentHistory{From: from, To: to, Amount: amount, Value: 1}}).
		Return(nil, nil)

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
	mockRates.EXPECT().
		GetLastRate(ctx, &currency.CurrencyTitle{Title: from}).
		Return(&currency.Currency{Title: from, Value: value}, nil)
	mockRates.EXPECT().
		GetLastRate(ctx, &currency.CurrencyTitle{Title: to}).
		Return(&currency.Currency{Title: to, Value: value}, nil)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestSetTransaction_FailDecode(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"

	req := httptest.NewRequest(http.MethodPost, url, nil)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSetTransactionFailDecode(t)
	defer ctrl.Finish()

	createContext(&req)
	testAuth.SetTransaction(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createSetTransactionFailDecode(t *testing.T) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestSetTransaction_FailCheckFrom400(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(fmt.Sprintf(`{"from": "%s", "to": "%s", "amount": %f}`, from, to, amount))

	req := httptest.NewRequest(http.MethodPost, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSetTransactionFailCheckFrom400(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.SetTransaction(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createSetTransactionFailCheckFrom400(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckWallet(ctx, createWalletCheck(from)).
		Return(&profileService.Check{Existence: false}, nil)

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestSetTransaction_FailCheckFrom500(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(fmt.Sprintf(`{"from": "%s", "to": "%s", "amount": %f}`, from, to, amount))

	req := httptest.NewRequest(http.MethodPost, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSetTransactionFailCheckFrom500(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.SetTransaction(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createSetTransactionFailCheckFrom500(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckWallet(ctx, createWalletCheck(from)).
		Return(&profileService.Check{Existence: false}, errors.New("createSetTransactionFailCheckFrom500"))

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestSetTransaction_FailCurrencyDivFrom(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(fmt.Sprintf(`{"from": "%s", "to": "%s", "amount": %f}`, from, to, amount))

	req := httptest.NewRequest(http.MethodPost, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSetTransactionFailCurrencyDivFrom(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.SetTransaction(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createSetTransactionFailCurrencyDivFrom(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckWallet(ctx, createWalletCheck(from)).
		Return(&profileService.Check{Existence: true}, nil)

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
	mockRates.EXPECT().
		GetLastRate(ctx, &currency.CurrencyTitle{Title: from}).
		Return(nil, errors.New("createSetTransactionFailCurrencyDivFrom"))

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestSetTransaction_FailCurrencyDivTo(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(fmt.Sprintf(`{"from": "%s", "to": "%s", "amount": %f}`, from, to, amount))

	req := httptest.NewRequest(http.MethodPost, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSetTransactionFailCurrencyDivTo(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.SetTransaction(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createSetTransactionFailCurrencyDivTo(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckWallet(ctx, createWalletCheck(from)).
		Return(&profileService.Check{Existence: true}, nil)

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
	mockRates.EXPECT().
		GetLastRate(ctx, &currency.CurrencyTitle{Title: from}).
		Return(&currency.Currency{Title: from, Value: value}, nil)
	mockRates.EXPECT().
		GetLastRate(ctx, &currency.CurrencyTitle{Title: to}).
		Return(nil, errors.New("createSetTransactionFailCurrencyDivTo"))

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestSetTransaction_FailGetWallet(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(fmt.Sprintf(`{"from": "%s", "to": "%s", "amount": %f}`, from, to, amount))

	req := httptest.NewRequest(http.MethodPost, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSetTransactionFailGetWallet(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.SetTransaction(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createSetTransactionFailGetWallet(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckWallet(ctx, createWalletCheck(from)).
		Return(&profileService.Check{Existence: true}, nil)
	mockUser.EXPECT().
		GetWallet(ctx, &profileService.ConcreteWallet{Id: id, Title: to}).
		Return(nil, errors.New("createSetTransactionFailGetWallet"))

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
	mockRates.EXPECT().
		GetLastRate(ctx, &currency.CurrencyTitle{Title: from}).
		Return(&currency.Currency{Title: from, Value: value}, nil)
	mockRates.EXPECT().
		GetLastRate(ctx, &currency.CurrencyTitle{Title: to}).
		Return(&currency.Currency{Title: to, Value: value}, nil)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestSetTransaction_FailCheckWalletTo(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(fmt.Sprintf(`{"from": "%s", "to": "%s", "amount": %f}`, from, to, amount))

	req := httptest.NewRequest(http.MethodPost, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSetTransactionFailCheckWalletTo(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.SetTransaction(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createSetTransactionFailCheckWalletTo(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckWallet(ctx, createWalletCheck(from)).
		Return(&profileService.Check{Existence: true}, nil)
	mockUser.EXPECT().
		GetWallet(ctx, &profileService.ConcreteWallet{Id: id, Title: to}).
		Return(&profileService.Wallet{Title: to, Value: walletSize}, nil)
	mockUser.EXPECT().
		CheckWallet(ctx, createWalletCheck(to)).
		Return(nil, errors.New("createSetTransactionFailCheckWalletTo"))

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
	mockRates.EXPECT().
		GetLastRate(ctx, &currency.CurrencyTitle{Title: from}).
		Return(&currency.Currency{Title: from, Value: value}, nil)
	mockRates.EXPECT().
		GetLastRate(ctx, &currency.CurrencyTitle{Title: to}).
		Return(&currency.Currency{Title: to, Value: value}, nil)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestSetTransaction_FailCheckWalletToCreateWallet(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(fmt.Sprintf(`{"from": "%s", "to": "%s", "amount": %f}`, from, to, amount))

	req := httptest.NewRequest(http.MethodPost, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSetTransactionFailCheckWalletToCreateWallet(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.SetTransaction(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createSetTransactionFailCheckWalletToCreateWallet(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckWallet(ctx, createWalletCheck(from)).
		Return(&profileService.Check{Existence: true}, nil)
	mockUser.EXPECT().
		GetWallet(ctx, &profileService.ConcreteWallet{Id: id, Title: to}).
		Return(&profileService.Wallet{Title: to, Value: walletSize}, nil)
	mockUser.EXPECT().
		CheckWallet(ctx, createWalletCheck(to)).
		Return(&profileService.Check{Existence: false}, nil)
	mockUser.EXPECT().
		CreateWallet(ctx, createWalletCheck(to)).
		Return(nil, errors.New("createSetTransactionFailCheckWalletToCreateWallet"))


	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
	mockRates.EXPECT().
		GetLastRate(ctx, &currency.CurrencyTitle{Title: from}).
		Return(&currency.Currency{Title: from, Value: value}, nil)
	mockRates.EXPECT().
		GetLastRate(ctx, &currency.CurrencyTitle{Title: to}).
		Return(&currency.Currency{Title: to, Value: value}, nil)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestSetTransaction_FailUpdateFrom(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(fmt.Sprintf(`{"from": "%s", "to": "%s", "amount": %f}`, from, to, amount))

	req := httptest.NewRequest(http.MethodPost, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSetTransactionFailUpdateFrom(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.SetTransaction(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createSetTransactionFailUpdateFrom(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckWallet(ctx, createWalletCheck(from)).
		Return(&profileService.Check{Existence: true}, nil)
	mockUser.EXPECT().
		GetWallet(ctx, &profileService.ConcreteWallet{Id: id, Title: to}).
		Return(&profileService.Wallet{Title: to, Value: walletSize}, nil)
	mockUser.EXPECT().
		CheckWallet(ctx, createWalletCheck(to)).
		Return(&profileService.Check{Existence: true}, nil)
	mockUser.EXPECT().
		UpdateWallet(ctx, createWalletToSet(from, -amount)).
		Return(nil, errors.New("createSetTransactionFailUpdateFrom"))

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
	mockRates.EXPECT().
		GetLastRate(ctx, &currency.CurrencyTitle{Title: from}).
		Return(&currency.Currency{Title: from, Value: value}, nil)
	mockRates.EXPECT().
		GetLastRate(ctx, &currency.CurrencyTitle{Title: to}).
		Return(&currency.Currency{Title: to, Value: value}, nil)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestSetTransaction_FailUpdateTo(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(fmt.Sprintf(`{"from": "%s", "to": "%s", "amount": %f}`, from, to, amount))

	req := httptest.NewRequest(http.MethodPost, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSetTransactionFailUpdateTo(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.SetTransaction(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createSetTransactionFailUpdateTo(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckWallet(ctx, createWalletCheck(from)).
		Return(&profileService.Check{Existence: true}, nil)
	mockUser.EXPECT().
		GetWallet(ctx, &profileService.ConcreteWallet{Id: id, Title: to}).
		Return(&profileService.Wallet{Title: to, Value: walletSize}, nil)
	mockUser.EXPECT().
		CheckWallet(ctx, createWalletCheck(to)).
		Return(&profileService.Check{Existence: true}, nil)
	mockUser.EXPECT().
		UpdateWallet(ctx, createWalletToSet(from, -amount)).
		Return(nil, nil)
	mockUser.EXPECT().
		UpdateWallet(ctx, createWalletToSet(to, amount)).
		Return(nil, errors.New("createSetTransactionFailUpdateTo"))

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
	mockRates.EXPECT().
		GetLastRate(ctx, &currency.CurrencyTitle{Title: from}).
		Return(&currency.Currency{Title: from, Value: value}, nil)
	mockRates.EXPECT().
		GetLastRate(ctx, &currency.CurrencyTitle{Title: to}).
		Return(&currency.Currency{Title: to, Value: value}, nil)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestSetTransaction_FailAddToPaymentHistory(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	body := strings.NewReader(fmt.Sprintf(`{"from": "%s", "to": "%s", "amount": %f}`, from, to, amount))

	req := httptest.NewRequest(http.MethodPost, url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSetTransactionFailAddToPaymentHistory(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.SetTransaction(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createSetTransactionFailAddToPaymentHistory(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckWallet(ctx, createWalletCheck(from)).
		Return(&profileService.Check{Existence: true}, nil)
	mockUser.EXPECT().
		GetWallet(ctx, &profileService.ConcreteWallet{Id: id, Title: to}).
		Return(&profileService.Wallet{Title: to, Value: walletSize}, nil)
	mockUser.EXPECT().
		CheckWallet(ctx, createWalletCheck(to)).
		Return(&profileService.Check{Existence: true}, nil)
	mockUser.EXPECT().
		UpdateWallet(ctx, createWalletToSet(from, -amount)).
		Return(nil, nil)
	mockUser.EXPECT().
		UpdateWallet(ctx, createWalletToSet(to, amount)).
		Return(nil, nil)
	mockUser.EXPECT().
		AddToPaymentHistory(ctx, &profileService.AddToHistory{Id: id, Transaction: &profileService.PaymentHistory{From: from, To: to, Amount: amount, Value: 1}}).
		Return(nil, errors.New("createSetTransactionFailAddToPaymentHistory"))

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
	mockRates.EXPECT().
		GetLastRate(ctx, &currency.CurrencyTitle{Title: from}).
		Return(&currency.Currency{Title: from, Value: value}, nil)
	mockRates.EXPECT().
		GetLastRate(ctx, &currency.CurrencyTitle{Title: to}).
		Return(&currency.Currency{Title: to, Value: value}, nil)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func createExpectedHistory() *profileService.AllHistory {
	return &profileService.AllHistory{History: []*profileService.PaymentHistory{{From: from, To: to, Amount: amount, Value: value}}}
}

func createWalletCheck(dest string) *profileService.ConcreteWallet {
	return &profileService.ConcreteWallet{Id: id, Title: dest}
}

func createWalletToSet(dest string, val float64) *profileService.ToSetWallet {
	return &profileService.ToSetWallet{Id: id, NewWallet: &profileService.Wallet{Title: dest, Value: val}}
}
