package mock

import (
	"github.com/golang/mock/gomock"
	"reflect"
	"server/src/domain/entity"
)

type DayCurrencyRepositoryForMock struct {
	ctrl     *gomock.Controller
	recorder *RecorderDayCurrencyMockRepository
}

type RecorderDayCurrencyMockRepository struct {
	mock *DayCurrencyRepositoryForMock
}

func NewDayCurrencyRepositoryForMock(ctrl *gomock.Controller) *DayCurrencyRepositoryForMock {
	mock := &DayCurrencyRepositoryForMock{ctrl: ctrl}
	mock.recorder = &RecorderDayCurrencyMockRepository{mock: mock}
	return mock
}

func (d *DayCurrencyRepositoryForMock) EXPECT() *RecorderDayCurrencyMockRepository {
	return d.recorder
}

func (d *DayCurrencyRepositoryForMock) SaveCurrency(currencies []entity.Currency) error {
	d.ctrl.T.Helper()
	ret := d.ctrl.Call(d, "SaveCurrency", currencies)
	err, _ := ret[0].(error)
	return err
}

func (r *RecorderFinanceMockRepository) SaveCurrency(currencies interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"SaveCurrency",
		reflect.TypeOf((*DayCurrencyRepositoryForMock)(nil).SaveCurrency),
		currencies,
	)
}

func (d *DayCurrencyRepositoryForMock) GetInitialCurrency() ([]entity.Currency, error) {
	d.ctrl.T.Helper()
	ret := d.ctrl.Call(d, "GetInitialCurrency")
	currencies, _ := ret[0].([]entity.Currency)
	err, _ := ret[1].(error)
	return currencies, err
}

func (r *RecorderFinanceMockRepository) GetInitialCurrency() *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"GetInitialCurrency",
		reflect.TypeOf((*DayCurrencyRepositoryForMock)(nil).GetInitialCurrency),
	)
}
