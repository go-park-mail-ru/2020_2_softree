package rates

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"server/src/canal/pkg/application"
	"server/src/canal/pkg/domain/entity"
	mocks "server/src/canal/pkg/infrastructure/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetRates_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/rates"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	testRate, ctrl := createForexRateSuccess(t)
	defer ctrl.Finish()

	testRate.GetRates(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func TestGetRates_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/rates"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	testRate, ctrl := createForexRateFail(t)
	defer ctrl.Finish()

	testRate.GetRates(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func TestRates_GetURLRateFail(t *testing.T) {
	url := "http://127.0.0.1:8000/api/rates/USD/"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	testRate, ctrl := createGetURLRateFail(t)
	defer ctrl.Finish()

	testRate.GetURLRate(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
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

func createForexRateSuccess(t *testing.T) (*Rates, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	rateMock := mocks.NewRateRepositoryForMock(ctrl)
	rateMock.EXPECT().GetRates().Return(createRates(), nil)

	dayCurrMock := mocks.NewDayCurrencyRepositoryForMock(ctrl)

	servicesDB := application.NewRateApp(rateMock, dayCurrMock)

	return NewRates(*servicesDB, servicesLog), ctrl
}

func createForexRateFail(t *testing.T) (*Rates, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	rateMock := mocks.NewRateRepositoryForMock(ctrl)
	rateMock.EXPECT().GetRates().Return(createRates(), errors.New("get rates"))

	dayCurrMock := mocks.NewDayCurrencyRepositoryForMock(ctrl)

	servicesDB := application.NewRateApp(rateMock, dayCurrMock)

	return NewRates(*servicesDB, servicesLog), ctrl
}

func createGetURLRateFail(t *testing.T) (*Rates, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	rateMock := mocks.NewRateRepositoryForMock(ctrl)

	dayCurrMock := mocks.NewDayCurrencyRepositoryForMock(ctrl)

	servicesDB := application.NewRateApp(rateMock, dayCurrMock)

	return NewRates(*servicesDB, servicesLog), ctrl
}

func createGetMarketsSuccess(t *testing.T) (*Rates, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	rateMock := mocks.NewRateRepositoryForMock(ctrl)

	dayCurrMock := mocks.NewDayCurrencyRepositoryForMock(ctrl)

	servicesDB := application.NewRateApp(rateMock, dayCurrMock)

	return NewRates(*servicesDB, servicesLog), ctrl
}

func createRates() []entity.Currency {
	base := "USD"
	currency := [...]string{"EUR", "RUB"}
	values := [...]float64{1.10, 0.23}

	rates := make([]entity.Currency, 0)
	rates = append(rates, entity.Currency{Base: base, Title: currency[0], Value: values[0]})
	rates = append(rates, entity.Currency{Base: base, Title: currency[1], Value: values[1]})

	return rates
}
