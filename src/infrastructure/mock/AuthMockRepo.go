package mock

import (
	"github.com/golang/mock/gomock"
	"reflect"
	"server/src/infrastructure/auth"
)

type AuthRepositoryForMock struct {
	ctrl     *gomock.Controller
	recorder *RecorderAuthMockRepository
}

type RecorderAuthMockRepository struct {
	mock *AuthRepositoryForMock
}

func NewAuthRepositoryForMock(ctrl *gomock.Controller) *AuthRepositoryForMock {
	mock := &AuthRepositoryForMock{ctrl: ctrl}
	mock.recorder = &RecorderAuthMockRepository{mock: mock}
	return mock
}

func (a *AuthRepositoryForMock) EXPECT() *RecorderAuthMockRepository {
	return a.recorder
}

func (a *AuthRepositoryForMock) CreateAuth(id uint64, val string) error {
	a.ctrl.T.Helper()
	ret := a.ctrl.Call(a, "CreateAuth", id, val)
	err, _ := ret[0].(error)
	return err
}

func (r *RecorderAuthMockRepository) CreateAuth(id, val interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"CreateAuth",
		reflect.TypeOf((*AuthRepositoryForMock)(nil).CreateAuth),
		id,
		val,
	)
}

func (a *AuthRepositoryForMock) CheckAuth(val string) (uint64, error) {
	a.ctrl.T.Helper()
	ret := a.ctrl.Call(a, "CheckAuth", val)
	id, _ := ret[0].(uint64)
	err, _ := ret[1].(error)
	return id, err
}

func (r *RecorderAuthMockRepository) CheckAuth(val interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"CheckAuth",
		reflect.TypeOf((*AuthRepositoryForMock)(nil).CheckAuth),
		val,
	)
}

func (a *AuthRepositoryForMock) DeleteAuth(details *auth.AccessDetails) error {
	a.ctrl.T.Helper()
	ret := a.ctrl.Call(a, "DeleteAuth", details)
	err, _ := ret[0].(error)
	return err
}

func (r *RecorderAuthMockRepository) DeleteAuth(details interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"DeleteAuth",
		reflect.TypeOf((*AuthRepositoryForMock)(nil).DeleteAuth),
		details,
	)
}
