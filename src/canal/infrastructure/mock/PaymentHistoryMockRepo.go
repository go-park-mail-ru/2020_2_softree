package mock

import (
	"github.com/golang/mock/gomock"
	"reflect"
	"server/src/canal/domain/entity"
)

type PaymentHistoryRepositoryForMock struct {
	ctrl     *gomock.Controller
	recorder *RecorderPaymentHistoryMockRepository
}

type RecorderPaymentHistoryMockRepository struct {
	mock *PaymentHistoryRepositoryForMock
}

func NewPaymentHistoryRepositoryForMock(ctrl *gomock.Controller) *PaymentHistoryRepositoryForMock {
	mock := &PaymentHistoryRepositoryForMock{ctrl: ctrl}
	mock.recorder = &RecorderPaymentHistoryMockRepository{mock: mock}
	return mock
}

func (p *PaymentHistoryRepositoryForMock) EXPECT() *RecorderPaymentHistoryMockRepository {
	return p.recorder
}

func (p *PaymentHistoryRepositoryForMock) GetAllPaymentHistory(id int64) ([]entity.PaymentHistory, error) {
	p.ctrl.T.Helper()
	ret := p.ctrl.Call(p, "GetAllPaymentHistory", id)
	history, _ := ret[0].([]entity.PaymentHistory)
	err, _ := ret[1].(error)
	return history, err
}

func (r *RecorderPaymentHistoryMockRepository) GetAllPaymentHistory(id interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"GetAllPaymentHistory",
		reflect.TypeOf((*PaymentHistoryRepositoryForMock)(nil).GetAllPaymentHistory),
		id,
	)
}

func (p *PaymentHistoryRepositoryForMock) GetIntervalPaymentHistory(id int64,
	i entity.Interval) ([]entity.PaymentHistory, error) {
	p.ctrl.T.Helper()
	ret := p.ctrl.Call(p, "GetIntervalPaymentHistory", id, i)
	history, _ := ret[0].([]entity.PaymentHistory)
	err, _ := ret[1].(error)
	return history, err
}

func (r *RecorderPaymentHistoryMockRepository) GetIntervalPaymentHistory(id, i interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"GetIntervalPaymentHistory",
		reflect.TypeOf((*PaymentHistoryRepositoryForMock)(nil).GetIntervalPaymentHistory),
		id,
		i,
	)
}

func (p *PaymentHistoryRepositoryForMock) AddToPaymentHistory(id int64, history entity.PaymentHistory) error {
	p.ctrl.T.Helper()
	ret := p.ctrl.Call(p, "AddToPaymentHistory", id, history)
	err, _ := ret[0].(error)
	return err
}

func (r *RecorderPaymentHistoryMockRepository) AddToPaymentHistory(id, history interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"AddToPaymentHistory",
		reflect.TypeOf((*PaymentHistoryRepositoryForMock)(nil).AddToPaymentHistory),
		id,
		history,
	)
}
