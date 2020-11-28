package mock

import (
	"github.com/golang/mock/gomock"
	"reflect"
	"server/src/currency/pkg/domain"
)

type ApiMock struct {
	ctrl     *gomock.Controller
	recorder *RecorderApi
}

type RecorderApi struct {
	mock *ApiMock
}

func NewApiMock(ctrl *gomock.Controller) *ApiMock {
	mock := &ApiMock{ctrl: ctrl}
	mock.recorder = &RecorderApi{mock: mock}
	return mock
}

func (api *ApiMock) EXPECT() *RecorderApi {
	return api.recorder
}

func (api *ApiMock) GetCurrencies() domain.FinancialRepository {
	api.ctrl.T.Helper()
	ret := api.ctrl.Call(api, "GetCurrencies")
	out, _ := ret[0].(domain.FinancialRepository)
	return out
}

func (r *RecorderApi) GetCurrencies() *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"GetCurrencies",
		reflect.TypeOf((*ApiMock)(nil).GetCurrencies),
	)
}
