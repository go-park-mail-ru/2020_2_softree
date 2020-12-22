package profile_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/infrastructure/mock"
	profileHTTP "server/canal/pkg/interfaces/profile"
	"testing"
)

func TestGetIncome_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	req = mux.SetURLVars(req, map[string]string{"period": period})
	testAuth, ctrl := createGetIncomeSuccess(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.GetIncome(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createGetIncomeSuccess(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)
	paymentLogic.EXPECT().
		GetIncome(ctx, entity.Income{Id: id, Period: period}).
		Return(entity.Description{}, decimal.New(1, 3), nil)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func TestGetIncome_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users"
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	req = mux.SetURLVars(req, map[string]string{"period": "week"})
	testAuth, ctrl := createGetIncomeFail(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.GetIncome(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createGetIncomeFail(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)
	paymentLogic.EXPECT().
		GetIncome(ctx, entity.Income{Id: id, Period: "week"}).
		Return(entity.Description{Status: 500}, decimal.New(1, 3), errors.New("error"))

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func TestGetIncomePerDay_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users?period=" + period
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	testAuth, ctrl := createGetIncomePerDaySuccess(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.GetAllIncomePerDay(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
}

func createGetIncomePerDaySuccess(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)
	paymentLogic.EXPECT().
		GetAllIncomePerDay(ctx, entity.Income{Id: id, Period: period}).
		Return(entity.Description{}, createAllIncome(), nil)

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func TestGetIncomePerDay_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/api/users?period=" + period
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	testAuth, ctrl := createGetIncomePerDayFail(t, createContext(&req))
	defer ctrl.Finish()

	testAuth.GetAllIncomePerDay(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
}

func createGetIncomePerDayFail(t *testing.T, ctx context.Context) (*profileHTTP.Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileLogic := mock.NewMockProfileLogic(ctrl)

	paymentLogic := mock.NewMockPaymentLogic(ctrl)
	paymentLogic.EXPECT().
		GetAllIncomePerDay(ctx, entity.Income{Id: id, Period: period}).
		Return(entity.Description{Status: 500}, entity.WalletStates{}, errors.New("error"))

	return profileHTTP.NewProfile(profileLogic, paymentLogic), ctrl
}

func createAllIncome() entity.WalletStates {
	return entity.WalletStates{{Value: decimal.NewFromFloat(value), UpdatedAt: ptypes.TimestampNow()}}
}
