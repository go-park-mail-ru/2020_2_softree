// Code generated by MockGen. DO NOT EDIT.
// Source: canal/pkg/domain/repository/payment.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	entity "server/canal/pkg/domain/entity"
)

// MockPaymentLogic is a mock of PaymentLogic interface
type MockPaymentLogic struct {
	ctrl     *gomock.Controller
	recorder *MockPaymentLogicMockRecorder
}

// MockPaymentLogicMockRecorder is the mock recorder for MockPaymentLogic
type MockPaymentLogicMockRecorder struct {
	mock *MockPaymentLogic
}

// NewMockPaymentLogic creates a new mock instance
func NewMockPaymentLogic(ctrl *gomock.Controller) *MockPaymentLogic {
	mock := &MockPaymentLogic{ctrl: ctrl}
	mock.recorder = &MockPaymentLogicMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPaymentLogic) EXPECT() *MockPaymentLogicMockRecorder {
	return m.recorder
}

// ReceiveTransactions mocks base method
func (m *MockPaymentLogic) ReceiveTransactions(ctx context.Context, id int64) (entity.Description, entity.Payments, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReceiveTransactions", ctx, id)
	ret0, _ := ret[0].(entity.Description)
	ret1, _ := ret[1].(entity.Payments)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ReceiveTransactions indicates an expected call of ReceiveTransactions
func (mr *MockPaymentLogicMockRecorder) ReceiveTransactions(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReceiveTransactions", reflect.TypeOf((*MockPaymentLogic)(nil).ReceiveTransactions), ctx, id)
}

// ReceiveWallets mocks base method
func (m *MockPaymentLogic) ReceiveWallets(ctx context.Context, id int64) (entity.Description, entity.Wallets, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReceiveWallets", ctx, id)
	ret0, _ := ret[0].(entity.Description)
	ret1, _ := ret[1].(entity.Wallets)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ReceiveWallets indicates an expected call of ReceiveWallets
func (mr *MockPaymentLogicMockRecorder) ReceiveWallets(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReceiveWallets", reflect.TypeOf((*MockPaymentLogic)(nil).ReceiveWallets), ctx, id)
}

// SetTransaction mocks base method
func (m *MockPaymentLogic) SetTransaction(ctx context.Context, payment entity.Payment) (entity.Description, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetTransaction", ctx, payment)
	ret0, _ := ret[0].(entity.Description)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetTransaction indicates an expected call of SetTransaction
func (mr *MockPaymentLogicMockRecorder) SetTransaction(ctx, payment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTransaction", reflect.TypeOf((*MockPaymentLogic)(nil).SetTransaction), ctx, payment)
}

// SetWallet mocks base method
func (m *MockPaymentLogic) SetWallet(ctx context.Context, wallet entity.Wallet) (entity.Description, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetWallet", ctx, wallet)
	ret0, _ := ret[0].(entity.Description)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetWallet indicates an expected call of SetWallet
func (mr *MockPaymentLogicMockRecorder) SetWallet(ctx, wallet interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWallet", reflect.TypeOf((*MockPaymentLogic)(nil).SetWallet), ctx, wallet)
}

// GetIncome mocks base method
func (m *MockPaymentLogic) GetIncome(ctx context.Context, in entity.Income) (entity.Description, float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIncome", ctx, in)
	ret0, _ := ret[0].(entity.Description)
	ret1, _ := ret[1].(float64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetIncome indicates an expected call of GetIncome
func (mr *MockPaymentLogicMockRecorder) GetIncome(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIncome", reflect.TypeOf((*MockPaymentLogic)(nil).GetIncome), ctx, in)
}

// WritePortfolios mocks base method
func (m *MockPaymentLogic) WritePortfolios() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WritePortfolios")
}

// WritePortfolios indicates an expected call of WritePortfolios
func (mr *MockPaymentLogicMockRecorder) WritePortfolios() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WritePortfolios", reflect.TypeOf((*MockPaymentLogic)(nil).WritePortfolios))
}
