// Code generated by MockGen. DO NOT EDIT.
// Source: currency.pb.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	reflect "reflect"
	gen "server/currency/pkg/currency/gen"
)

// MockCurrencyServiceClient is a mock of CurrencyServiceClient interface
type MockCurrencyServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockCurrencyServiceClientMockRecorder
}

// MockCurrencyServiceClientMockRecorder is the mock recorder for MockCurrencyServiceClient
type MockCurrencyServiceClientMockRecorder struct {
	mock *MockCurrencyServiceClient
}

// NewMockCurrencyServiceClient creates a new mock instance
func NewMockCurrencyServiceClient(ctrl *gomock.Controller) *MockCurrencyServiceClient {
	mock := &MockCurrencyServiceClient{ctrl: ctrl}
	mock.recorder = &MockCurrencyServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCurrencyServiceClient) EXPECT() *MockCurrencyServiceClientMockRecorder {
	return m.recorder
}

// GetAllLatestRates mocks base method
func (m *MockCurrencyServiceClient) GetAllLatestRates(ctx context.Context, in *gen.Empty, opts ...grpc.CallOption) (*gen.Currencies, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAllLatestRates", varargs...)
	ret0, _ := ret[0].(*gen.Currencies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllLatestRates indicates an expected call of GetAllLatestRates
func (mr *MockCurrencyServiceClientMockRecorder) GetAllLatestRates(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllLatestRates", reflect.TypeOf((*MockCurrencyServiceClient)(nil).GetAllLatestRates), varargs...)
}

// GetAllRatesByTitle mocks base method
func (m *MockCurrencyServiceClient) GetAllRatesByTitle(ctx context.Context, in *gen.CurrencyTitle, opts ...grpc.CallOption) (*gen.Currencies, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAllRatesByTitle", varargs...)
	ret0, _ := ret[0].(*gen.Currencies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllRatesByTitle indicates an expected call of GetAllRatesByTitle
func (mr *MockCurrencyServiceClientMockRecorder) GetAllRatesByTitle(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllRatesByTitle", reflect.TypeOf((*MockCurrencyServiceClient)(nil).GetAllRatesByTitle), varargs...)
}

// GetLastRate mocks base method
func (m *MockCurrencyServiceClient) GetLastRate(ctx context.Context, in *gen.CurrencyTitle, opts ...grpc.CallOption) (*gen.Currency, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetLastRate", varargs...)
	ret0, _ := ret[0].(*gen.Currency)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastRate indicates an expected call of GetLastRate
func (mr *MockCurrencyServiceClientMockRecorder) GetLastRate(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastRate", reflect.TypeOf((*MockCurrencyServiceClient)(nil).GetLastRate), varargs...)
}

// GetInitialDayCurrency mocks base method
func (m *MockCurrencyServiceClient) GetInitialDayCurrency(ctx context.Context, in *gen.Empty, opts ...grpc.CallOption) (*gen.InitialDayCurrencies, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetInitialDayCurrency", varargs...)
	ret0, _ := ret[0].(*gen.InitialDayCurrencies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInitialDayCurrency indicates an expected call of GetInitialDayCurrency
func (mr *MockCurrencyServiceClientMockRecorder) GetInitialDayCurrency(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInitialDayCurrency", reflect.TypeOf((*MockCurrencyServiceClient)(nil).GetInitialDayCurrency), varargs...)
}

// MockCurrencyServiceServer is a mock of CurrencyServiceServer interface
type MockCurrencyServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockCurrencyServiceServerMockRecorder
}

// MockCurrencyServiceServerMockRecorder is the mock recorder for MockCurrencyServiceServer
type MockCurrencyServiceServerMockRecorder struct {
	mock *MockCurrencyServiceServer
}

// NewMockCurrencyServiceServer creates a new mock instance
func NewMockCurrencyServiceServer(ctrl *gomock.Controller) *MockCurrencyServiceServer {
	mock := &MockCurrencyServiceServer{ctrl: ctrl}
	mock.recorder = &MockCurrencyServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCurrencyServiceServer) EXPECT() *MockCurrencyServiceServerMockRecorder {
	return m.recorder
}

// GetAllLatestRates mocks base method
func (m *MockCurrencyServiceServer) GetAllLatestRates(arg0 context.Context, arg1 *gen.Empty) (*gen.Currencies, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllLatestRates", arg0, arg1)
	ret0, _ := ret[0].(*gen.Currencies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllLatestRates indicates an expected call of GetAllLatestRates
func (mr *MockCurrencyServiceServerMockRecorder) GetAllLatestRates(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllLatestRates", reflect.TypeOf((*MockCurrencyServiceServer)(nil).GetAllLatestRates), arg0, arg1)
}

// GetAllRatesByTitle mocks base method
func (m *MockCurrencyServiceServer) GetAllRatesByTitle(arg0 context.Context, arg1 *gen.CurrencyTitle) (*gen.Currencies, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllRatesByTitle", arg0, arg1)
	ret0, _ := ret[0].(*gen.Currencies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllRatesByTitle indicates an expected call of GetAllRatesByTitle
func (mr *MockCurrencyServiceServerMockRecorder) GetAllRatesByTitle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllRatesByTitle", reflect.TypeOf((*MockCurrencyServiceServer)(nil).GetAllRatesByTitle), arg0, arg1)
}

// GetLastRate mocks base method
func (m *MockCurrencyServiceServer) GetLastRate(arg0 context.Context, arg1 *gen.CurrencyTitle) (*gen.Currency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastRate", arg0, arg1)
	ret0, _ := ret[0].(*gen.Currency)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastRate indicates an expected call of GetLastRate
func (mr *MockCurrencyServiceServerMockRecorder) GetLastRate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastRate", reflect.TypeOf((*MockCurrencyServiceServer)(nil).GetLastRate), arg0, arg1)
}

// GetInitialDayCurrency mocks base method
func (m *MockCurrencyServiceServer) GetInitialDayCurrency(arg0 context.Context, arg1 *gen.Empty) (*gen.InitialDayCurrencies, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInitialDayCurrency", arg0, arg1)
	ret0, _ := ret[0].(*gen.InitialDayCurrencies)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInitialDayCurrency indicates an expected call of GetInitialDayCurrency
func (mr *MockCurrencyServiceServerMockRecorder) GetInitialDayCurrency(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInitialDayCurrency", reflect.TypeOf((*MockCurrencyServiceServer)(nil).GetInitialDayCurrency), arg0, arg1)
}
