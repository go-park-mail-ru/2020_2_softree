package mock

import (
	"context"
	"github.com/golang/mock/gomock"
	"reflect"
	profile "server/src/profile/pkg/profile/gen"
)

func (profileMock *ProfileMock) GetWallet(ctx context.Context, in *profile.ConcreteWallet) (*profile.Wallet, error) {
	profileMock.ctrl.T.Helper()
	ret := profileMock.ctrl.Call(profileMock, "GetWallet", ctx, in)
	out, _ := ret[0].(*profile.Wallet)
	err, _ := ret[1].(error)
	return out, err
}

func (profileRecorder *RecorderProfile) GetWallet(ctx, in interface{}) *gomock.Call {
	profileRecorder.mock.ctrl.T.Helper()
	return profileRecorder.mock.ctrl.RecordCallWithMethodType(
		profileRecorder.mock,
		"GetWallet",
		reflect.TypeOf((*ProfileMock)(nil).GetWallet),
		ctx,
		in,
	)
}

func (profileMock *ProfileMock) SetWallet(ctx context.Context, in *profile.ToSetWallet) (*profile.Empty, error) {
	profileMock.ctrl.T.Helper()
	ret := profileMock.ctrl.Call(profileMock, "SetWallet", ctx, in)
	out, _ := ret[0].(*profile.Empty)
	err, _ := ret[1].(error)
	return out, err
}

func (profileRecorder *RecorderProfile) SetWallet(ctx, in interface{}) *gomock.Call {
	profileRecorder.mock.ctrl.T.Helper()
	return profileRecorder.mock.ctrl.RecordCallWithMethodType(
		profileRecorder.mock,
		"SetWallet",
		reflect.TypeOf((*ProfileMock)(nil).SetWallet),
		ctx,
		in,
	)
}

func (profileMock *ProfileMock) GetWallets(ctx context.Context, in *profile.UserID) (*profile.Wallets, error) {
	profileMock.ctrl.T.Helper()
	ret := profileMock.ctrl.Call(profileMock, "GetWallets", ctx, in)
	out, _ := ret[0].(*profile.Wallets)
	err, _ := ret[1].(error)
	return out, err
}

func (profileRecorder *RecorderProfile) GetWallets(ctx, in interface{}) *gomock.Call {
	profileRecorder.mock.ctrl.T.Helper()
	return profileRecorder.mock.ctrl.RecordCallWithMethodType(
		profileRecorder.mock,
		"GetWallets",
		reflect.TypeOf((*ProfileMock)(nil).GetWallets),
		ctx,
		in,
	)
}

func (profileMock *ProfileMock) CreateWallet(ctx context.Context, in *profile.ConcreteWallet) (*profile.Empty, error) {
	profileMock.ctrl.T.Helper()
	ret := profileMock.ctrl.Call(profileMock, "CreateWallet", ctx, in)
	out, _ := ret[0].(*profile.Empty)
	err, _ := ret[1].(error)
	return out, err
}

func (profileRecorder *RecorderProfile) CreateWallet(ctx, in interface{}) *gomock.Call {
	profileRecorder.mock.ctrl.T.Helper()
	return profileRecorder.mock.ctrl.RecordCallWithMethodType(
		profileRecorder.mock,
		"CreateWallet",
		reflect.TypeOf((*ProfileMock)(nil).CreateWallet),
		ctx,
		in,
	)
}

func (profileMock *ProfileMock) CheckWallet(ctx context.Context, in *profile.ConcreteWallet) (*profile.Check, error) {
	profileMock.ctrl.T.Helper()
	ret := profileMock.ctrl.Call(profileMock, "CheckWallet", ctx, in)
	out, _ := ret[0].(*profile.Check)
	err, _ := ret[1].(error)
	return out, err
}

func (profileRecorder *RecorderProfile) CheckWallet(ctx, in interface{}) *gomock.Call {
	profileRecorder.mock.ctrl.T.Helper()
	return profileRecorder.mock.ctrl.RecordCallWithMethodType(
		profileRecorder.mock,
		"CheckWallet",
		reflect.TypeOf((*ProfileMock)(nil).CheckWallet),
		ctx,
		in,
	)
}

func (profileMock *ProfileMock) UpdateWallet(ctx context.Context, in *profile.ToSetWallet) (*profile.Empty, error) {
	profileMock.ctrl.T.Helper()
	ret := profileMock.ctrl.Call(profileMock, "UpdateWallet", ctx, in)
	out, _ := ret[0].(*profile.Empty)
	err, _ := ret[1].(error)
	return out, err
}

func (profileRecorder *RecorderProfile) UpdateWallet(ctx, in interface{}) *gomock.Call {
	profileRecorder.mock.ctrl.T.Helper()
	return profileRecorder.mock.ctrl.RecordCallWithMethodType(
		profileRecorder.mock,
		"UpdateWallet",
		reflect.TypeOf((*ProfileMock)(nil).UpdateWallet),
		ctx,
		in,
	)
}
