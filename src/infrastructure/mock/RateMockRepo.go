package mock

import (
	"github.com/golang/mock/gomock"
	"reflect"
	"server/src/domain/entity"
	"server/src/domain/repository"
)

type RateRepositoryForMock struct {
	ctrl     *gomock.Controller
	recorder *RecorderRateMockRepository
}

type RecorderRateMockRepository struct {
	mock *RateRepositoryForMock
}

func NewRateRepositoryForMock(ctrl *gomock.Controller) *RateRepositoryForMock {
	mock := &RateRepositoryForMock{ctrl: ctrl}
	mock.recorder = &RecorderRateMockRepository{mock: mock}
	return mock
}

func (a *RateRepositoryForMock) EXPECT() *RecorderRateMockRepository {
	return a.recorder
}

func (a *RateRepositoryForMock) SaveRates(financial repository.FinancialRepository) error {
	a.ctrl.T.Helper()
	ret := a.ctrl.Call(a, "SaveRates", financial)
	err, _ := ret[0].(error)
	return err
}

func (r *RecorderRateMockRepository) SaveRates(financial interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"SaveRates",
		reflect.TypeOf((*RateRepositoryForMock)(nil).SaveRates),
		financial,
	)
}

func (a *RateRepositoryForMock) UpdateRate(u uint64, rate entity.Currency) (entity.Currency, error) {
	panic("implement me")
}

func (a *RateRepositoryForMock) DeleteRate(u uint64) error {
	panic("implement me")
}

func (a *RateRepositoryForMock) GetRates() ([]entity.Currency, error) {
	a.ctrl.T.Helper()
	ret := a.ctrl.Call(a, "GetRates")
	rates, _ := ret[0].([]entity.Currency)
	err, _ := ret[1].(error)
	return rates, err
}

func (r *RecorderRateMockRepository) GetRates() *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"GetRates",
		reflect.TypeOf((*RateRepositoryForMock)(nil).GetRates),
	)
}

func (a *RateRepositoryForMock) GetRate(title string) ([]entity.Currency, error) {
	a.ctrl.T.Helper()
	ret := a.ctrl.Call(a, "GetRate", title)
	rates, _ := ret[0].([]entity.Currency)
	err, _ := ret[1].(error)
	return rates, err
}

func (r *RecorderRateMockRepository) GetRate(title interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"GetRate",
		reflect.TypeOf((*RateRepositoryForMock)(nil).GetRate),
		title,
	)
}
