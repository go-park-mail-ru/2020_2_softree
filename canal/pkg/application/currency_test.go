package application_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"reflect"
	"server/canal/pkg/application"
	"server/canal/pkg/domain/entity"
	"server/currency/pkg/currency/gen"
	currency "server/currency/pkg/infrastructure/mock"
	"testing"
)

func TestGetAllLatestCurrencies_Success(t *testing.T) {
	req := createRequestNotInitial()
	testAuth, ctrl := createGetAllLatestCurrenciesSuccess(t, req.Context())
	defer ctrl.Finish()

	desc, out, err := testAuth.GetAllLatestCurrencies(req)

	require.NoError(t, err)
	require.Empty(t, desc)
	require.NotEmpty(t, out)
	require.Equal(t, reflect.TypeOf(entity.Currencies{}), reflect.TypeOf(out))
}

func createGetAllLatestCurrenciesSuccess(t *testing.T, ctx context.Context) (*application.CurrencyApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	currencyService.EXPECT().
		GetAllLatestRates(ctx, &gen.Empty{}).
		Return(createCurrencies(), nil)

	return application.NewCurrencyApp(currencyService), ctrl
}

func TestGetAllLatestCurrencies_Fail(t *testing.T) {
	req := createRequestNotInitial()
	testAuth, ctrl := createGetAllLatestCurrenciesFail(t, req.Context())
	defer ctrl.Finish()

	desc, out, err := testAuth.GetAllLatestCurrencies(req)

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Empty(t, out)
	require.Equal(t, reflect.TypeOf(entity.Currencies{}), reflect.TypeOf(out))
}

func createGetAllLatestCurrenciesFail(t *testing.T, ctx context.Context) (*application.CurrencyApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	currencyService.EXPECT().
		GetAllLatestRates(ctx, &gen.Empty{}).
		Return(&gen.Currencies{}, errors.New("error"))

	return application.NewCurrencyApp(currencyService), ctrl
}

func TestGetURLCurrencies_Success(t *testing.T) {
	req := createRequestNotInitial()
	testAuth, ctrl := createGetURLCurrenciesSuccess(t, req.Context())
	defer ctrl.Finish()

	desc, out, err := testAuth.GetURLCurrencies(req)

	require.NoError(t, err)
	require.Empty(t, desc)
	require.NotEmpty(t, out)
	require.Equal(t, reflect.TypeOf(entity.Currencies{}), reflect.TypeOf(out))
}

func createGetURLCurrenciesSuccess(t *testing.T, ctx context.Context) (*application.CurrencyApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	currencyService.EXPECT().
		GetAllRatesByTitle(ctx, &gen.CurrencyTitle{Title: curr}).
		Return(createCurrencies(), nil)

	return application.NewCurrencyApp(currencyService), ctrl
}

func TestGetAllLatestCurrenciesInitial_Success(t *testing.T) {
	req := createRequestInitial()
	testAuth, ctrl := createGetAllLatestCurrenciesInitialSuccess(t, req.Context())
	defer ctrl.Finish()

	desc, out, err := testAuth.GetAllLatestCurrencies(req)

	require.NoError(t, err)
	require.Empty(t, desc)
	require.NotEmpty(t, out)
	require.Equal(t, reflect.TypeOf(entity.Currencies{}), reflect.TypeOf(out))
}

func createGetAllLatestCurrenciesInitialSuccess(t *testing.T, ctx context.Context) (*application.CurrencyApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	currencyService.EXPECT().
		GetInitialDayCurrency(ctx, &gen.Empty{}).
		Return(createInitialCurrencies(), nil)

	return application.NewCurrencyApp(currencyService), ctrl
}

func TestGetAllLatestCurrenciesInitial_Fail(t *testing.T) {
	req := createRequestInitial()
	testAuth, ctrl := createGetAllLatestCurrenciesInitialFail(t, req.Context())
	defer ctrl.Finish()

	desc, out, err := testAuth.GetAllLatestCurrencies(req)

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Empty(t, out)
	require.Equal(t, reflect.TypeOf(entity.Currencies{}), reflect.TypeOf(out))
}

func createGetAllLatestCurrenciesInitialFail(t *testing.T, ctx context.Context) (*application.CurrencyApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	currencyService.EXPECT().
		GetInitialDayCurrency(ctx, &gen.Empty{}).
		Return(&gen.InitialDayCurrencies{}, errors.New("error"))

	return application.NewCurrencyApp(currencyService), ctrl
}

func TestGetURLCurrencies_Fail(t *testing.T) {
	ctx := createRequestNotInitial()
	testAuth, ctrl := createGetURLCurrenciesFail(t, ctx.Context())
	defer ctrl.Finish()

	desc, out, err := testAuth.GetURLCurrencies(ctx)

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Empty(t, out)
	require.Equal(t, reflect.TypeOf(entity.Currencies{}), reflect.TypeOf(out))
}

func createGetURLCurrenciesFail(t *testing.T, ctx context.Context) (*application.CurrencyApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	currencyService.EXPECT().
		GetAllRatesByTitle(ctx, &gen.CurrencyTitle{Title: curr}).
		Return(&gen.Currencies{}, errors.New("error"))

	return application.NewCurrencyApp(currencyService), ctrl
}

func TestGetURLCurrencies_FailValidate(t *testing.T) {
	ctx := createFailRequest()
	testAuth, ctrl := createGetURLCurrenciesFailValidate(t, ctx.Context())
	defer ctrl.Finish()

	desc, out, err := testAuth.GetURLCurrencies(ctx)

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Empty(t, out)
	require.Equal(t, reflect.TypeOf(entity.Currencies{}), reflect.TypeOf(out))
}

func createGetURLCurrenciesFailValidate(t *testing.T, ctx context.Context) (*application.CurrencyApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)

	return application.NewCurrencyApp(currencyService), ctrl
}

func TestGetMarkets_Success(t *testing.T) {
	req := createRequestNotInitial()
	testAuth, ctrl := createGetMarketsSuccess(t, req.Context())
	defer ctrl.Finish()

	desc, out, err := testAuth.GetMarkets()

	require.NoError(t, err)
	require.Empty(t, desc)
	require.NotEmpty(t, out)
	require.Equal(t, reflect.TypeOf(entity.Markets{}), reflect.TypeOf(out))
}

func createGetMarketsSuccess(t *testing.T, ctx context.Context) (*application.CurrencyApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)

	return application.NewCurrencyApp(currencyService), ctrl
}

func createCurrencies() *gen.Currencies {
	return &gen.Currencies{Rates: []*gen.Currency{{Base: "USD", Title: "RUB", Value: 0.23, UpdatedAt: ptypes.TimestampNow()}}}
}

func createInitialCurrencies() *gen.InitialDayCurrencies {
	return &gen.InitialDayCurrencies{Currencies: []*gen.InitialDayCurrency{{Title: "RUB", Value: 0.23}}}
}

func createRequestNotInitial() *http.Request {
	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8000/api/rates/USD", nil)
	req = mux.SetURLVars(req, map[string]string{"title": curr})
	return req
}

func createRequestInitial() *http.Request {
	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8000/api/rates/USD?initial=true", nil)
	req = mux.SetURLVars(req, map[string]string{"title": curr})
	return req
}

func createFailRequest() *http.Request {
	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8000/api/rates/some", nil)
	req = mux.SetURLVars(req, map[string]string{"title": "some"})
	return req
}
