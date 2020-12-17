package application_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"reflect"
	"server/canal/pkg/application"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/infrastructure/mock"
	profile "server/profile/pkg/infrastructure/mock"
	"server/profile/pkg/profile/gen"
	"testing"
)

func TestReceiveTransactions_Success(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createReceiveTransactionsSuccess(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.ReceiveTransactions(ctx, id)

	require.NoError(t, err)
	require.Empty(t, desc)
	require.NotEmpty(t, out)
	require.Equal(t, reflect.TypeOf(entity.Payments{}), reflect.TypeOf(out))
}

func createReceiveTransactionsSuccess(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		GetAllPaymentHistory(ctx, &gen.UserID{Id: id}).
		Return(createHistory(), nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, securityService), ctrl
}

func TestReceiveTransactions_Fail(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createReceiveTransactionsFail(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.ReceiveTransactions(ctx, id)

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Empty(t, out)
	require.Equal(t, reflect.TypeOf(entity.Payments{}), reflect.TypeOf(out))
}

func createReceiveTransactionsFail(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		GetAllPaymentHistory(ctx, &gen.UserID{Id: id}).
		Return(&gen.AllHistory{}, errors.New("error"))

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, securityService), ctrl
}

func createHistory() *gen.AllHistory {
	return &gen.AllHistory{
		History: []*gen.PaymentHistory{
			{Base: base, Currency: curr, Value: value, Amount: amount},
		},
	}
}
