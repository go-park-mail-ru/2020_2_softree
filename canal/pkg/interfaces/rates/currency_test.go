package rates

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	currencyService "server/currency/pkg/currency/gen"
	"server/currency/pkg/infrastructure/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const title = "USD"

func TestGetRates_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/rates"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	testRate, ctrl := createForexRateSuccess(t, req.Context())
	defer ctrl.Finish()

	testRate.GetAllLatestRates(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createForexRateSuccess(t *testing.T, ctx context.Context) (*Rates, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	rateMock := mock.NewRateRepositoryForMock(ctrl)
	rateMock.EXPECT().
		GetAllLatestRates(ctx, &currencyService.Empty{}).
		Return(createRates(), nil)

	return NewRates(rateMock), ctrl
}

func TestGetRates_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/rates"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	testRate, ctrl := createForexRateFail(t, req.Context())
	defer ctrl.Finish()

	testRate.GetAllLatestRates(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func createForexRateFail(t *testing.T, ctx context.Context) (*Rates, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	rateMock := mock.NewRateRepositoryForMock(ctrl)
	rateMock.EXPECT().
		GetAllLatestRates(ctx, &currencyService.Empty{}).
		Return(nil, errors.New("createForexRateFail"))

	return NewRates(rateMock), ctrl
}

func TestGetURLRate_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/api/rates/USD"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	req = mux.SetURLVars(req, map[string]string{"title": title})
	testRate, ctrl := createGetURLRateSuccess(t, req.Context())
	defer ctrl.Finish()

	testRate.GetURLRate(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createGetURLRateSuccess(t *testing.T, ctx context.Context) (*Rates, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	rateMock := mock.NewRateRepositoryForMock(ctrl)
	rateMock.EXPECT().
		GetAllRatesByTitle(ctx, &currencyService.CurrencyTitle{Title: title}).
		Return(createRates(), nil)

	return NewRates(rateMock), ctrl
}

func TestGetURLRate_FailGetRate(t *testing.T) {
	url := "http://127.0.0.1:8000/api/rates/USD"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	req = mux.SetURLVars(req, map[string]string{"title": title})
	testRate, ctrl := createGetURLRateFailGetRate(t, req.Context())
	defer ctrl.Finish()

	testRate.GetURLRate(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createGetURLRateFailGetRate(t *testing.T, ctx context.Context) (*Rates, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	rateMock := mock.NewRateRepositoryForMock(ctrl)
	rateMock.EXPECT().
		GetAllRatesByTitle(ctx, &currencyService.CurrencyTitle{Title: title}).
		Return(nil, errors.New("createGetURLRateFailGetRate"))

	return NewRates(rateMock), ctrl
}

func TestRates_GetURLRateFail(t *testing.T) {
	url := "http://127.0.0.1:8000/api/rates/USD"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	testRate, ctrl := createGetURLRateFail(t)
	defer ctrl.Finish()

	testRate.GetURLRate(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createGetURLRateFail(t *testing.T) (*Rates, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	rateMock := mock.NewRateRepositoryForMock(ctrl)

	return NewRates(rateMock), ctrl
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

func createGetMarketsSuccess(t *testing.T) (*Rates, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	rateMock := mock.NewRateRepositoryForMock(ctrl)

	return NewRates(rateMock), ctrl
}

func createRates() *currencyService.Currencies {
	base := "USD"
	currency := [...]string{"EUR", "RUB"}
	values := [...]float64{1.10, 0.23}

	var rates currencyService.Currencies
	rates.Rates = append(rates.Rates, &currencyService.Currency{Base: base, Title: currency[0], Value: values[0]})
	rates.Rates = append(rates.Rates, &currencyService.Currency{Base: base, Title: currency[1], Value: values[1]})

	return &rates
}
