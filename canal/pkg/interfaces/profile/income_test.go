package profile_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"server/canal/pkg/infrastructure/mock"
	profileHTTP "server/canal/pkg/interfaces/profile"
	currencyService "server/currency/pkg/currency/gen"
	currencyMock "server/currency/pkg/infrastructure/mock"
	profileMock "server/profile/pkg/infrastructure/mock"
	profileService "server/profile/pkg/profile/gen"
	"testing"
)

func TestGetIncome_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	req = mux.SetURLVars(req, map[string]string{"period": "day"})
	testAuth, ctrl := createGetIncomeSuccess(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.GetIncome(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createGetIncomeSuccess(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		GetIncome(ctx, &profileService.IncomeParameters{Id: id, Period: "day"}).
		Return(&profileService.Income{Change: 200}, nil)
	mockUser.EXPECT().
		GetWallets(ctx, &profileService.UserID{Id: id}).
		Return(createExpectedWallets(), nil)

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
	mockRates.EXPECT().
		GetLastRate(ctx, &currencyService.CurrencyTitle{Title: curr}).
		Return(&currencyService.Currency{Title: curr, Value: 1}, nil)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestGetIncome_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	req = mux.SetURLVars(req, map[string]string{"period": "day"})
	testAuth, ctrl := createGetIncomeFail(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.GetIncome(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createGetIncomeFail(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		GetIncome(ctx, &profileService.IncomeParameters{Id: id, Period: "day"}).
		Return(&profileService.Income{Change: 0}, errors.New("createGetIncomeFail"))

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestGetIncome_FailGetWallets(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	req = mux.SetURLVars(req, map[string]string{"period": "day"})
	testAuth, ctrl := createGetIncomeFailGetWallets(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.GetIncome(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createGetIncomeFailGetWallets(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		GetIncome(ctx, &profileService.IncomeParameters{Id: id, Period: "day"}).
		Return(&profileService.Income{Change: 200}, nil)
	mockUser.EXPECT().
		GetWallets(ctx, &profileService.UserID{Id: id}).
		Return(createExpectedWallets(), errors.New("createGetIncomeFailGetWallets"))

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}

func TestGetIncome_FailGetLastRate(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	req = mux.SetURLVars(req, map[string]string{"period": "day"})
	testAuth, ctrl := createGetIncomeFailGetLastRate(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.GetIncome(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createGetIncomeFailGetLastRate(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		GetIncome(ctx, &profileService.IncomeParameters{Id: id, Period: "day"}).
		Return(&profileService.Income{Change: 200}, nil)
	mockUser.EXPECT().
		GetWallets(ctx, &profileService.UserID{Id: id}).
		Return(createExpectedWallets(), nil)

	mockSecurity := mock.NewSecurityMock(ctrl)

	mockRates := currencyMock.NewRateRepositoryForMock(ctrl)
	mockRates.EXPECT().
		GetLastRate(ctx, &currencyService.CurrencyTitle{Title: curr}).
		Return(&currencyService.Currency{Title: curr, Value: 1}, errors.New("createGetIncomeFailGetLastRate"))

	return profileHTTP.NewProfile(mockUser, mockSecurity, mockRates), ctrl
}
