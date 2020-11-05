package rates

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"server/src/application"
	"server/src/domain/entity"
	"server/src/domain/repository"
	"server/src/infrastructure/log"
	mocks "server/src/infrastructure/mock"
	"testing"
)

func TestGetRates_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/rates"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	testRate, repo, ctrl := createForexRateSuccess(t)
	defer ctrl.Finish()

	ctx := context.WithValue(req.Context(), "finance", repo)
	req = req.Clone(ctx)

	testRate.GetRates(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func TestGetRates_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/rates"

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	testRate, repo, ctrl := createForexRateFail(t)
	defer ctrl.Finish()

	ctx := context.WithValue(req.Context(), "finance", repo)
	req = req.Clone(ctx)

	testRate.GetRates(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createForexRateSuccess(t *testing.T) (*Rates, repository.FinancialRepository, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	financeMock := mocks.NewFinanceRepositoryForMock(ctrl)

	rateMock := mocks.NewRateRepositoryForMock(ctrl)
	rateMock.EXPECT().SaveRates(financeMock).Return(createRates(), nil)

	servicesDB := application.NewRateApp(rateMock)
	servicesLog := log.NewLogrusLogger()

	return NewRates(*servicesDB, servicesLog), financeMock, ctrl
}

func createForexRateFail(t *testing.T) (*Rates, repository.FinancialRepository, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	financeMock := mocks.NewFinanceRepositoryForMock(ctrl)

	rateMock := mocks.NewRateRepositoryForMock(ctrl)
	rateMock.EXPECT().SaveRates(financeMock).Return(createRates(), errors.New("error"))

	servicesDB := application.NewRateApp(rateMock)
	servicesLog := log.NewLogrusLogger()

	return NewRates(*servicesDB, servicesLog), financeMock, ctrl
}

func createFinanceMap() map[string]interface{} {
	financeMap := make(map[string]interface{}, 1)
	financeMap["EUR"] = 1.10
	financeMap["RUB"] = 0.23

	return financeMap
}

func createRates() []entity.Rate {
	base := "USD"
	currency := [...]string{"EUR", "RUB"}
	values := [...]float64{1.10, 0.23}

	rates := make([]entity.Rate, 0)
	rates = append(rates, entity.Rate{Base: base, Currency: currency[0], Value: values[0], ID: 1})
	rates = append(rates, entity.Rate{Base: base, Currency: currency[1], Value: values[1], ID: 2})

	return rates
}
