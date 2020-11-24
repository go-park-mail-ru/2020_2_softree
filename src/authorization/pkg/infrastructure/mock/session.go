package mock

import (
	"context"
	"github.com/golang/mock/gomock"
	"reflect"
	session "server/src/authorization/pkg/session/gen"
)

type AuthMock struct {
	ctrl     *gomock.Controller
	recorder *RecorderAuth
}

type RecorderAuth struct {
	mock *AuthMock
}

func NewAuthRepositoryForMock(ctrl *gomock.Controller) *AuthMock {
	mock := &AuthMock{ctrl: ctrl}
	mock.recorder = &RecorderAuth{mock: mock}
	return mock
}

func (auth *AuthMock) EXPECT() *RecorderAuth {
	return auth.recorder
}

func (auth *AuthMock) Create(ctx context.Context, in *session.Session) (*session.UserID, error) {
	auth.ctrl.T.Helper()
	ret := auth.ctrl.Call(auth, "Create", ctx, in)
	out, _ := ret[0].(*session.UserID)
	err, _ := ret[1].(error)
	return out, err
}

func (recorderAuth *RecorderAuth) Create(ctx, in interface{}) *gomock.Call {
	recorderAuth.mock.ctrl.T.Helper()
	return recorderAuth.mock.ctrl.RecordCallWithMethodType(
		recorderAuth.mock,
		"Create",
		reflect.TypeOf((*AuthMock)(nil).Create),
		ctx,
		in,
	)
}

func (auth *AuthMock) Check(ctx context.Context, in *session.SessionID) (*session.UserID, error) {
	auth.ctrl.T.Helper()
	ret := auth.ctrl.Call(auth, "Check", ctx, in)
	out, _ := ret[0].(*session.UserID)
	err, _ := ret[1].(error)
	return out, err
}

func (recorderAuth *RecorderAuth) Check(ctx, in interface{}) *gomock.Call {
	recorderAuth.mock.ctrl.T.Helper()
	return recorderAuth.mock.ctrl.RecordCallWithMethodType(
		recorderAuth.mock,
		"Check",
		reflect.TypeOf((*AuthMock)(nil).Check),
		ctx,
		in,
	)
}

func (auth *AuthMock) Delete(ctx context.Context, in *session.SessionID) (*session.Empty, error) {
	auth.ctrl.T.Helper()
	ret := auth.ctrl.Call(auth, "Delete", ctx, in)
	out, _ := ret[0].(*session.Empty)
	err, _ := ret[1].(error)
	return out, err
}

func (recorderAuth *RecorderAuth) Delete(ctx, in interface{}) *gomock.Call {
	recorderAuth.mock.ctrl.T.Helper()
	return recorderAuth.mock.ctrl.RecordCallWithMethodType(
		recorderAuth.mock,
		"Delete",
		reflect.TypeOf((*AuthMock)(nil).Delete),
		ctx,
		in,
	)
}
