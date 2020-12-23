// Code generated by MockGen. DO NOT EDIT.
// Source: financialAPI.go

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	domain "server/currency/pkg/domain"
)

// MockFinancialAPI is a mock of FinancialAPI interface
type MockFinancialAPI struct {
	ctrl     *gomock.Controller
	recorder *MockFinancialAPIMockRecorder
}

// MockFinancialAPIMockRecorder is the mock recorder for MockFinancialAPI
type MockFinancialAPIMockRecorder struct {
	mock *MockFinancialAPI
}

// NewMockFinancialAPI creates a new mock instance
func NewMockFinancialAPI(ctrl *gomock.Controller) *MockFinancialAPI {
	mock := &MockFinancialAPI{ctrl: ctrl}
	mock.recorder = &MockFinancialAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFinancialAPI) EXPECT() *MockFinancialAPIMockRecorder {
	return m.recorder
}

// GetCurrencies mocks base method
func (m *MockFinancialAPI) GetCurrencies() (domain.FinancialRepository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrencies")
	ret0, _ := ret[0].(domain.FinancialRepository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrencies indicates an expected call of GetCurrencies
func (mr *MockFinancialAPIMockRecorder) GetCurrencies() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrencies", reflect.TypeOf((*MockFinancialAPI)(nil).GetCurrencies))
}
