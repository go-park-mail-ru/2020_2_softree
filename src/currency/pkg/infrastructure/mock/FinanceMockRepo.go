package mock

import (
	"github.com/golang/mock/gomock"
	"reflect"
)

type FinanceRepositoryForMock struct {
	ctrl     *gomock.Controller
	recorder *RecorderFinanceMockRepository
}

type RecorderFinanceMockRepository struct {
	mock *FinanceRepositoryForMock
}

func NewFinanceRepositoryForMock(ctrl *gomock.Controller) *FinanceRepositoryForMock {
	mock := &FinanceRepositoryForMock{ctrl: ctrl}
	mock.recorder = &RecorderFinanceMockRepository{mock: mock}
	return mock
}

func (f *FinanceRepositoryForMock) EXPECT() *RecorderFinanceMockRepository {
	return f.recorder
}

func (f *FinanceRepositoryForMock) GetQuote() map[string]interface{} {
	f.ctrl.T.Helper()
	ret := f.ctrl.Call(f, "GetQuote")
	quote, _ := ret[0].(map[string]interface{})
	return quote
}

func (r *RecorderFinanceMockRepository) GetQuote() *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"GetQuote",
		reflect.TypeOf((*FinanceRepositoryForMock)(nil).GetQuote),
	)
}

func (f *FinanceRepositoryForMock) GetBase() string {
	f.ctrl.T.Helper()
	ret := f.ctrl.Call(f, "GetBase")
	base, _ := ret[0].(string)
	return base
}

func (r *RecorderFinanceMockRepository) GetBase() *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"GetBase",
		reflect.TypeOf((*FinanceRepositoryForMock)(nil).GetBase),
	)
}
