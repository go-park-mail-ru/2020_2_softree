package rates_test

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"net/http"
	"net/http/httptest"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/infrastructure/mock"
	"server/canal/pkg/interfaces/rates"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAllLatestRates_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/rates"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	testRate, ctrl := createGetAllLatestRatesSuccess(t, req)
	defer ctrl.Finish()

	testRate.GetAllLatestRates(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createGetAllLatestRatesSuccess(t *testing.T, req *http.Request) (*rates.Rates, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	currencyLogic := mock.NewMockCurrencyLogic(ctrl)
	currencyLogic.EXPECT().
		GetAllLatestCurrencies(req).
		Return(entity.Description{}, createCurrencies(), nil)

	return rates.NewRates(currencyLogic), ctrl
}

func TestGetAllLatestRates_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/rates"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	testRate, ctrl := createGetAllLatestRatesFail(t, req)
	defer ctrl.Finish()

	testRate.GetAllLatestRates(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func createGetAllLatestRatesFail(t *testing.T, req *http.Request) (*rates.Rates, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	currencyLogic := mock.NewMockCurrencyLogic(ctrl)
	currencyLogic.EXPECT().
		GetAllLatestCurrencies(req).
		Return(entity.Description{Status: 500}, entity.Currencies{}, errors.New("error"))

	return rates.NewRates(currencyLogic), ctrl
}

func TestGetURLRate_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/api/rates/USD"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	req = mux.SetURLVars(req, map[string]string{"title": "USD"})
	testRate, ctrl := createGetURLRateSuccess(t, req)
	defer ctrl.Finish()

	testRate.GetURLRate(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createGetURLRateSuccess(t *testing.T, r *http.Request) (*rates.Rates, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	currencyLogic := mock.NewMockCurrencyLogic(ctrl)
	currencyLogic.EXPECT().
		GetURLCurrencies(r).
		Return(entity.Description{}, createCurrencies(), nil)

	return rates.NewRates(currencyLogic), ctrl
}

func TestGetURLRate_FailGetRate(t *testing.T) {
	url := "http://127.0.0.1:8000/api/rates/USD"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	req = mux.SetURLVars(req, map[string]string{"title": "USD"})
	testRate, ctrl := createGetURLRateFailGetRate(t, req)
	defer ctrl.Finish()

	testRate.GetURLRate(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createGetURLRateFailGetRate(t *testing.T, r *http.Request) (*rates.Rates, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	currencyLogic := mock.NewMockCurrencyLogic(ctrl)
	currencyLogic.EXPECT().
		GetURLCurrencies(r).
		Return(entity.Description{Status: 500}, entity.Currencies{}, errors.New("error"))

	return rates.NewRates(currencyLogic), ctrl
}

func TestRates_GetMarketsSuccess(t *testing.T) {
	url := "http://127.0.0.1:8000/api/markets"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	testRate, ctrl := createGetMarketsSuccess(t)
	defer ctrl.Finish()

	testRate.GetMarkets(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createGetMarketsSuccess(t *testing.T) (*rates.Rates, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	currencyLogic := mock.NewMockCurrencyLogic(ctrl)
	currencyLogic.EXPECT().
		GetMarkets().
		Return(entity.Description{}, createMarkets(), nil)

	return rates.NewRates(currencyLogic), ctrl
}

func TestRates_GetMarketsFail(t *testing.T) {
	url := "http://127.0.0.1:8000/api/markets"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	testRate, ctrl := createGetMarketsFail(t)
	defer ctrl.Finish()

	testRate.GetMarkets(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createGetMarketsFail(t *testing.T) (*rates.Rates, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	currencyLogic := mock.NewMockCurrencyLogic(ctrl)
	currencyLogic.EXPECT().
		GetMarkets().
		Return(entity.Description{Status: 500}, entity.Markets{}, errors.New("error"))

	return rates.NewRates(currencyLogic), ctrl
}

func createCurrencies() entity.Currencies {
	return entity.Currencies{{Base: "USD", Title: "RUB", Value: decimal.NewFromFloat(0.23)}}
}

func createMarkets() entity.Markets {
	return entity.Markets{{Base: "USD", Title: "RUB"}}
}
