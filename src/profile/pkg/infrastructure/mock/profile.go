package mock

import (
	"context"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
	"reflect"
	profile "server/src/profile/pkg/profile/gen"
)

type ProfileMock struct {
	ctrl     *gomock.Controller
	recorder *RecorderProfile
}

type RecorderProfile struct {
	mock *ProfileMock
}

func NewProfileMock(ctrl *gomock.Controller) *ProfileMock {
	mock := &ProfileMock{ctrl: ctrl}
	mock.recorder = &RecorderProfile{mock: mock}
	return mock
}

func (profileMock *ProfileMock) EXPECT() *RecorderProfile {
	return profileMock.recorder
}

func (profileMock *ProfileMock) SaveUser(ctx context.Context, in *profile.User, opts ...grpc.CallOption) (*profile.PublicUser, error) {
	profileMock.ctrl.T.Helper()
	ret := profileMock.ctrl.Call(profileMock, "SaveUser", ctx, in)
	out, _ := ret[0].(*profile.PublicUser)
	err, _ := ret[1].(error)
	return out, err
}

func (profileRecorder *RecorderProfile) SaveUser(ctx, in interface{}) *gomock.Call {
	profileRecorder.mock.ctrl.T.Helper()
	return profileRecorder.mock.ctrl.RecordCallWithMethodType(
		profileRecorder.mock,
		"SaveUser",
		reflect.TypeOf((*ProfileMock)(nil).SaveUser),
		ctx,
		in,
	)
}

func (profileMock *ProfileMock) UpdateUserAvatar(ctx context.Context, in *profile.UpdateFields, opts ...grpc.CallOption) (*profile.Empty, error) {
	profileMock.ctrl.T.Helper()
	ret := profileMock.ctrl.Call(profileMock, "UpdateUserAvatar", ctx, in)
	out, _ := ret[0].(*profile.Empty)
	err, _ := ret[1].(error)
	return out, err
}

func (profileRecorder *RecorderProfile) UpdateUserAvatar(ctx, in interface{}) *gomock.Call {
	profileRecorder.mock.ctrl.T.Helper()
	return profileRecorder.mock.ctrl.RecordCallWithMethodType(
		profileRecorder.mock,
		"UpdateUserAvatar",
		reflect.TypeOf((*ProfileMock)(nil).UpdateUserAvatar),
		ctx,
		in,
	)
}

func (profileMock *ProfileMock) UpdateUserPassword(ctx context.Context, in *profile.UpdateFields, opts ...grpc.CallOption) (*profile.Empty, error) {
	profileMock.ctrl.T.Helper()
	ret := profileMock.ctrl.Call(profileMock, "UpdateUserPassword", ctx, in)
	out, _ := ret[0].(*profile.Empty)
	err, _ := ret[1].(error)
	return out, err
}

func (profileRecorder *RecorderProfile) UpdateUserPassword(ctx, in interface{}) *gomock.Call {
	profileRecorder.mock.ctrl.T.Helper()
	return profileRecorder.mock.ctrl.RecordCallWithMethodType(
		profileRecorder.mock,
		"UpdateUserPassword",
		reflect.TypeOf((*ProfileMock)(nil).UpdateUserPassword),
		ctx,
		in,
	)
}

func (profileMock *ProfileMock) DeleteUser(ctx context.Context, in *profile.UserID, opts ...grpc.CallOption) (*profile.Empty, error) {
	profileMock.ctrl.T.Helper()
	ret := profileMock.ctrl.Call(profileMock, "DeleteUser", ctx, in)
	out, _ := ret[0].(*profile.Empty)
	err, _ := ret[1].(error)
	return out, err
}

func (profileRecorder *RecorderProfile) DeleteUser(ctx, in interface{}) *gomock.Call {
	profileRecorder.mock.ctrl.T.Helper()
	return profileRecorder.mock.ctrl.RecordCallWithMethodType(
		profileRecorder.mock,
		"DeleteUser",
		reflect.TypeOf((*ProfileMock)(nil).DeleteUser),
		ctx,
		in,
	)
}

func (profileMock *ProfileMock) GetUserById(ctx context.Context, in *profile.UserID, opts ...grpc.CallOption) (*profile.PublicUser, error) {
	profileMock.ctrl.T.Helper()
	ret := profileMock.ctrl.Call(profileMock, "GetUserById", ctx, in)
	out, _ := ret[0].(*profile.PublicUser)
	err, _ := ret[1].(error)
	return out, err
}

func (profileRecorder *RecorderProfile) GetUserById(ctx, in interface{}) *gomock.Call {
	profileRecorder.mock.ctrl.T.Helper()
	return profileRecorder.mock.ctrl.RecordCallWithMethodType(
		profileRecorder.mock,
		"GetUserById",
		reflect.TypeOf((*ProfileMock)(nil).GetUserById),
		ctx,
		in,
	)
}

func (profileMock *ProfileMock) GetUserByLogin(ctx context.Context, in *profile.User, opts ...grpc.CallOption) (*profile.PublicUser, error) {
	profileMock.ctrl.T.Helper()
	ret := profileMock.ctrl.Call(profileMock, "GetUserByLogin", ctx, in)
	out, _ := ret[0].(*profile.PublicUser)
	err, _ := ret[1].(error)
	return out, err
}

func (profileRecorder *RecorderProfile) GetUserByLogin(ctx, in interface{}) *gomock.Call {
	profileRecorder.mock.ctrl.T.Helper()
	return profileRecorder.mock.ctrl.RecordCallWithMethodType(
		profileRecorder.mock,
		"GetUserByLogin",
		reflect.TypeOf((*ProfileMock)(nil).GetUserByLogin),
		ctx,
		in,
	)
}

func (profileMock *ProfileMock) CheckExistence(ctx context.Context, in *profile.User, opts ...grpc.CallOption) (*profile.Check, error) {
	profileMock.ctrl.T.Helper()
	ret := profileMock.ctrl.Call(profileMock, "CheckExistence", ctx, in)
	out, _ := ret[0].(*profile.Check)
	err, _ := ret[1].(error)
	return out, err
}

func (profileRecorder *RecorderProfile) CheckExistence(ctx, in interface{}) *gomock.Call {
	profileRecorder.mock.ctrl.T.Helper()
	return profileRecorder.mock.ctrl.RecordCallWithMethodType(
		profileRecorder.mock,
		"CheckExistence",
		reflect.TypeOf((*ProfileMock)(nil).CheckExistence),
		ctx,
		in,
	)
}

func (profileMock *ProfileMock) CheckPassword(ctx context.Context, in *profile.User, opts ...grpc.CallOption) (*profile.Check, error) {
	profileMock.ctrl.T.Helper()
	ret := profileMock.ctrl.Call(profileMock, "CheckPassword", ctx, in)
	out, _ := ret[0].(*profile.Check)
	err, _ := ret[1].(error)
	return out, err
}

func (profileRecorder *RecorderProfile) CheckPassword(ctx, in interface{}) *gomock.Call {
	profileRecorder.mock.ctrl.T.Helper()
	return profileRecorder.mock.ctrl.RecordCallWithMethodType(
		profileRecorder.mock,
		"CheckPassword",
		reflect.TypeOf((*ProfileMock)(nil).CheckExistence),
		ctx,
		in,
	)
}

func (profileMock *ProfileMock) GetUserWatchlist(ctx context.Context, in *profile.UserID, opts ...grpc.CallOption) (*profile.Currencies, error) {
	profileMock.ctrl.T.Helper()
	ret := profileMock.ctrl.Call(profileMock, "GetUserWatchlist", ctx, in)
	out, _ := ret[0].(*profile.Currencies)
	err, _ := ret[1].(error)
	return out, err
}

func (profileRecorder *RecorderProfile) GetUserWatchlist(ctx, in interface{}) *gomock.Call {
	profileRecorder.mock.ctrl.T.Helper()
	return profileRecorder.mock.ctrl.RecordCallWithMethodType(
		profileRecorder.mock,
		"GetUserWatchlist",
		reflect.TypeOf((*ProfileMock)(nil).GetUserWatchlist),
		ctx,
		in,
	)
}
