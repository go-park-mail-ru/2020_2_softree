package mock

import (
	"context"
	"github.com/golang/mock/gomock"
	"reflect"
	session "server/src/authorization/pkg/session/gen"
)

type SessionMock struct {
	ctrl     *gomock.Controller
	recorder *RecorderSession
}

type RecorderSession struct {
	mock *SessionMock
}

func NewAuthRepositoryForMock(ctrl *gomock.Controller) *SessionMock {
	mock := &SessionMock{ctrl: ctrl}
	mock.recorder = &RecorderSession{mock: mock}
	return mock
}

func (sessionMock *SessionMock) EXPECT() *RecorderSession {
	return sessionMock.recorder
}

func (sessionMock *SessionMock) Create(ctx context.Context, in *session.Session) (*session.UserID, error) {
	sessionMock.ctrl.T.Helper()
	ret := sessionMock.ctrl.Call(sessionMock, "Create", ctx, in)
	out, _ := ret[0].(*session.UserID)
	err, _ := ret[1].(error)
	return out, err
}

func (recorderAuth *RecorderSession) Create(ctx, in interface{}) *gomock.Call {
	recorderAuth.mock.ctrl.T.Helper()
	return recorderAuth.mock.ctrl.RecordCallWithMethodType(
		recorderAuth.mock,
		"Create",
		reflect.TypeOf((*SessionMock)(nil).Create),
		ctx,
		in,
	)
}

func (sessionMock *SessionMock) Check(ctx context.Context, in *session.SessionID) (*session.UserID, error) {
	sessionMock.ctrl.T.Helper()
	ret := sessionMock.ctrl.Call(sessionMock, "Check", ctx, in)
	out, _ := ret[0].(*session.UserID)
	err, _ := ret[1].(error)
	return out, err
}

func (recorderAuth *RecorderSession) Check(ctx, in interface{}) *gomock.Call {
	recorderAuth.mock.ctrl.T.Helper()
	return recorderAuth.mock.ctrl.RecordCallWithMethodType(
		recorderAuth.mock,
		"Check",
		reflect.TypeOf((*SessionMock)(nil).Check),
		ctx,
		in,
	)
}

func (sessionMock *SessionMock) Delete(ctx context.Context, in *session.SessionID) (*session.Empty, error) {
	sessionMock.ctrl.T.Helper()
	ret := sessionMock.ctrl.Call(sessionMock, "Delete", ctx, in)
	out, _ := ret[0].(*session.Empty)
	err, _ := ret[1].(error)
	return out, err
}

func (recorderAuth *RecorderSession) Delete(ctx, in interface{}) *gomock.Call {
	recorderAuth.mock.ctrl.T.Helper()
	return recorderAuth.mock.ctrl.RecordCallWithMethodType(
		recorderAuth.mock,
		"Delete",
		reflect.TypeOf((*SessionMock)(nil).Delete),
		ctx,
		in,
	)
}
