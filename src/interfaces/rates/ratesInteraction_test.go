package rates

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"server/src/application"
	"server/src/domain/entity"
	"server/src/infrastructure/log"
	mocks "server/src/infrastructure/mock"
	"testing"
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

func createForexRateSuccess(t *testing.T) (*Rates, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	rateMock := mocks.NewRateRepositoryForMock(ctrl)
	rateMock.EXPECT().GetRates().Return(createRates(), nil)

	servicesDB := application.NewRateApp(rateMock)
	servicesLog := log.NewLogrusLogger()

	return NewRates(*servicesDB, servicesLog), ctrl
}

func createForexRateFail(t *testing.T) (*Rates, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	rateMock := mocks.NewRateRepositoryForMock(ctrl)
	rateMock.EXPECT().GetRates().Return(createRates(), errors.New("get rates"))

	servicesDB := application.NewRateApp(rateMock)
	servicesLog := log.NewLogrusLogger()

	return NewRates(*servicesDB, servicesLog), ctrl
}

func createRates() []entity.Currency {
	base := "USD"
	currency := [...]string{"EUR", "RUB"}
	values := [...]float64{1.10, 0.23}

	rates := make([]entity.Currency, 0)
	rates = append(rates, entity.Currency{Base: base, Title: currency[0], Value: values[0], ID: 1})
	rates = append(rates, entity.Currency{Base: base, Title: currency[1], Value: values[1], ID: 2})

	return rates
}
