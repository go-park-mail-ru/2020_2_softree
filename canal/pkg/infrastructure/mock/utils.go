package mock

import (
	"github.com/golang/mock/gomock"
	"reflect"
)

type SecurityMock struct {
	ctrl     *gomock.Controller
	recorder *RecorderSecurity
}

type RecorderSecurity struct {
	mock *SecurityMock
}

func NewSecurityMock(ctrl *gomock.Controller) *SecurityMock {
	mock := &SecurityMock{ctrl: ctrl}
	mock.recorder = &RecorderSecurity{mock: mock}
	return mock
}

func (d *SecurityMock) EXPECT() *RecorderSecurity {
	return d.recorder
}

func (d *SecurityMock) MakeShieldedPassword(stringToHash string) (string, error) {
	d.ctrl.T.Helper()
	ret := d.ctrl.Call(d, "MakeShieldedPassword", stringToHash)
	out, _ := ret[0].(string)
	err, _ := ret[1].(error)
	return out, err
}

func (r *RecorderSecurity) MakeShieldedPassword(stringToHash interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"MakeShieldedPassword",
		reflect.TypeOf((*SecurityMock)(nil).MakeShieldedPassword),
		stringToHash,
	)
}
