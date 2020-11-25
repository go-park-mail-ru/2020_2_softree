package mock

import (
	"context"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
	"reflect"
	profile "server/src/profile/pkg/profile/gen"
)

func (profileMock *ProfileMock) GetAllPaymentHistory(ctx context.Context, in *profile.UserID, opts ...grpc.CallOption) (*profile.AllHistory, error) {
	profileMock.ctrl.T.Helper()
	ret := profileMock.ctrl.Call(profileMock, "GetAllPaymentHistory", ctx, in)
	out, _ := ret[0].(*profile.AllHistory)
	err, _ := ret[1].(error)
	return out, err
}

func (profileRecorder *RecorderProfile) GetAllPaymentHistory(ctx, in interface{}) *gomock.Call {
	profileRecorder.mock.ctrl.T.Helper()
	return profileRecorder.mock.ctrl.RecordCallWithMethodType(
		profileRecorder.mock,
		"GetAllPaymentHistory",
		reflect.TypeOf((*ProfileMock)(nil).GetAllPaymentHistory),
		ctx,
		in,
	)
}

func (profileMock *ProfileMock) AddToPaymentHistory(ctx context.Context, in *profile.AddToHistory, opts ...grpc.CallOption) (*profile.Empty, error) {
	profileMock.ctrl.T.Helper()
	ret := profileMock.ctrl.Call(profileMock, "AddToPaymentHistory", ctx, in)
	out, _ := ret[0].(*profile.Empty)
	err, _ := ret[1].(error)
	return out, err
}

func (profileRecorder *RecorderProfile) AddToPaymentHistory(ctx, in interface{}) *gomock.Call {
	profileRecorder.mock.ctrl.T.Helper()
	return profileRecorder.mock.ctrl.RecordCallWithMethodType(
		profileRecorder.mock,
		"AddToPaymentHistory",
		reflect.TypeOf((*ProfileMock)(nil).AddToPaymentHistory),
		ctx,
		in,
	)
}
