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
	currencyGen "server/currency/pkg/currency/gen"
	currency "server/currency/pkg/infrastructure/mock"
	profile "server/profile/pkg/infrastructure/mock"
	profileGen "server/profile/pkg/profile/gen"
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
		GetAllPaymentHistory(ctx, &profileGen.UserID{Id: id}).
		Return(createHistory(), nil)

	currencyService := currency.NewRateRepositoryForMock(ctrl)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
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
		GetAllPaymentHistory(ctx, &profileGen.UserID{Id: id}).
		Return(&profileGen.AllHistory{}, errors.New("error"))

	currencyService := currency.NewRateRepositoryForMock(ctrl)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestReceiveWallets_Success(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createReceiveWalletsSuccess(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.ReceiveWallets(ctx, id)

	require.NoError(t, err)
	require.Empty(t, desc)
	require.NotEmpty(t, out)
	require.Equal(t, reflect.TypeOf(entity.Wallets{}), reflect.TypeOf(out))
}

func createReceiveWalletsSuccess(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		GetWallets(ctx, &profileGen.UserID{Id: id}).
		Return(createWallets(), nil)

	currencyService := currency.NewRateRepositoryForMock(ctrl)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestReceiveWallets_Fail(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createReceiveWalletsFail(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.ReceiveWallets(ctx, id)

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Empty(t, out)
	require.Equal(t, reflect.TypeOf(entity.Wallets{}), reflect.TypeOf(out))
}

func createReceiveWalletsFail(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		GetWallets(ctx, &profileGen.UserID{Id: id}).
		Return(&profileGen.Wallets{}, errors.New("error"))

	currencyService := currency.NewRateRepositoryForMock(ctrl)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestSetWallets_Success(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createSetWalletsSuccess(t, ctx)
	defer ctrl.Finish()

	desc, err := testAuth.SetWallet(ctx, entity.Wallet{Title: curr, UserId: id})

	require.NoError(t, err)
	require.Empty(t, desc)
}

func createSetWalletsSuccess(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		CreateWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Empty{}, nil)

	currencyService := currency.NewRateRepositoryForMock(ctrl)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestSetWallets_Fail(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createSetWalletsFail(t, ctx)
	defer ctrl.Finish()

	desc, err := testAuth.SetWallet(ctx, entity.Wallet{Title: curr, UserId: id})

	require.Error(t, err)
	require.NotEmpty(t, desc)
}

func createSetWalletsFail(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		CreateWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Empty{}, errors.New("error"))

	currencyService := currency.NewRateRepositoryForMock(ctrl)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestSetTransaction_Success(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createSetTransactionSuccess(t, ctx)
	defer ctrl.Finish()

	desc, err := testAuth.SetWallet(ctx, entity.Wallet{Title: curr, UserId: id})

	require.NoError(t, err)
	require.Empty(t, desc)
}

func createSetTransactionSuccess(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)

	currencyService := currency.NewRateRepositoryForMock(ctrl)
	currencyService.EXPECT().
		GetLastRate(ctx, currencyGen.CurrencyTitle{Title: base}).
		Return(baseValue, nil)
	currencyService.EXPECT().
		GetLastRate(ctx, currencyGen.CurrencyTitle{Title: curr}).
		Return(currValue, nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func createHistory() *profileGen.AllHistory {
	return &profileGen.AllHistory{
		History: []*profileGen.PaymentHistory{
			{Base: base, Currency: curr, Value: currValue, Amount: amount},
		},
	}
}

func createWallets() *profileGen.Wallets {
	return &profileGen.Wallets{
		Wallets: []*profileGen.Wallet{
			{Title: curr, Value: currValue},
		},
	}
}
