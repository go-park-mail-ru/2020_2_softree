package mock

import (
	"github.com/golang/mock/gomock"
	"net/http"
	"reflect"
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

func (a *AuthRepositoryForMock) CreateAuth(id int64) (http.Cookie, error) {
	a.ctrl.T.Helper()
	ret := a.ctrl.Call(a, "CreateAuth", id)
	cookie, _ := ret[0].(http.Cookie)
	err, _ := ret[1].(error)
	return cookie, err
}

func (r *RecorderAuthMockRepository) CreateAuth(id interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"CreateAuth",
		reflect.TypeOf((*AuthRepositoryForMock)(nil).CreateAuth),
		id,
	)
}

func (a *AuthRepositoryForMock) CheckAuth(val string) (int64, error) {
	a.ctrl.T.Helper()
	ret := a.ctrl.Call(a, "CheckAuth", val)
	id, _ := ret[0].(int64)
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

func (a *AuthRepositoryForMock) DeleteAuth(details string) error {
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

func (a *AuthRepositoryForMock) CreateCookie() (http.Cookie, error) {
	a.ctrl.T.Helper()
	ret := a.ctrl.Call(a, "CreateCookie")
	cookie, _ := ret[0].(http.Cookie)
	err, _ := ret[1].(error)
	return cookie, err
}

func (r *RecorderAuthMockRepository) CreateCookie() *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"CreateCookie",
		reflect.TypeOf((*AuthRepositoryForMock)(nil).CreateCookie),
	)
}
