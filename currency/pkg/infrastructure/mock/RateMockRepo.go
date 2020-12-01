package mock

import (
	"context"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
	"reflect"
	currency "server/currency/pkg/currency/gen"
)

type RateMock struct {
	ctrl     *gomock.Controller
	recorder *RecorderRate
}

type RecorderRate struct {
	mock *RateMock
}

func NewRateRepositoryForMock(ctrl *gomock.Controller) *RateMock {
	mock := &RateMock{ctrl: ctrl}
	mock.recorder = &RecorderRate{mock: mock}
	return mock
}

func (rateMock *RateMock) EXPECT() *RecorderRate {
	return rateMock.recorder
}

func (rateMock *RateMock) GetAllLatestRates(ctx context.Context, in *currency.Empty, opts ...grpc.CallOption) (*currency.Currencies, error) {
	rateMock.ctrl.T.Helper()
	ret := rateMock.ctrl.Call(rateMock, "GetAllLatestRates", ctx, in)
	out, _ := ret[0].(*currency.Currencies)
	err, _ := ret[1].(error)
	return out, err
}

func (recorder *RecorderRate) GetAllLatestRates(ctx, in interface{}) *gomock.Call {
	recorder.mock.ctrl.T.Helper()
	return recorder.mock.ctrl.RecordCallWithMethodType(
		recorder.mock,
		"GetAllLatestRates",
		reflect.TypeOf((*RateMock)(nil).GetAllLatestRates),
		ctx,
		in,
	)
}

func (rateMock *RateMock) GetAllRatesByTitle(ctx context.Context, in *currency.CurrencyTitle, opts ...grpc.CallOption) (*currency.Currencies, error) {
	rateMock.ctrl.T.Helper()
	ret := rateMock.ctrl.Call(rateMock, "GetAllRatesByTitle", ctx, in)
	out, _ := ret[0].(*currency.Currencies)
	err, _ := ret[1].(error)
	return out, err
}

func (recorder *RecorderRate) GetAllRatesByTitle(ctx, in interface{}) *gomock.Call {
	recorder.mock.ctrl.T.Helper()
	return recorder.mock.ctrl.RecordCallWithMethodType(
		recorder.mock,
		"GetAllRatesByTitle",
		reflect.TypeOf((*RateMock)(nil).GetAllRatesByTitle),
		ctx,
		in,
	)
}

func (rateMock *RateMock) GetInitialDayCurrency(ctx context.Context, in *currency.Empty, opts ...grpc.CallOption) (*currency.InitialDayCurrencies, error) {
	rateMock.ctrl.T.Helper()
	ret := rateMock.ctrl.Call(rateMock, "GetInitialDayCurrency")
	out, _ := ret[0].(*currency.InitialDayCurrencies)
	err, _ := ret[1].(error)
	return out, err
}

func (recorder *RecorderRate) GetInitialCurrency(ctx, in interface{}) *gomock.Call {
	recorder.mock.ctrl.T.Helper()
	return recorder.mock.ctrl.RecordCallWithMethodType(
		recorder.mock,
		"GetInitialDayCurrency",
		reflect.TypeOf((*RateMock)(nil).GetInitialDayCurrency),
		ctx,
		in,
	)
}

func (rateMock *RateMock) GetLastRate(ctx context.Context, in *currency.CurrencyTitle, opts ...grpc.CallOption) (*currency.Currency, error) {
	rateMock.ctrl.T.Helper()
	ret := rateMock.ctrl.Call(rateMock, "GetLastRate", ctx, in)
	out, _ := ret[0].(*currency.Currency)
	err, _ := ret[1].(error)
	return out, err
}

func (recorder *RecorderRate) GetLastRate(ctx, in interface{}) *gomock.Call {
	recorder.mock.ctrl.T.Helper()
	return recorder.mock.ctrl.RecordCallWithMethodType(
		recorder.mock,
		"GetLastRate",
		reflect.TypeOf((*RateMock)(nil).GetLastRate),
		ctx,
		in,
	)
}
