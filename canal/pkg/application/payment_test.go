package application_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
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

	profileService := profile.NewMockProfileServiceClient(ctrl)
	profileService.EXPECT().
		GetAllPaymentHistory(ctx, &profileGen.UserID{Id: id}).
		Return(createHistory(), nil)

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)

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

	profileService := profile.NewMockProfileServiceClient(ctrl)
	profileService.EXPECT().
		GetAllPaymentHistory(ctx, &profileGen.UserID{Id: id}).
		Return(&profileGen.AllHistory{}, errors.New("error"))

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)

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

	profileService := profile.NewMockProfileServiceClient(ctrl)
	profileService.EXPECT().
		GetWallets(ctx, &profileGen.UserID{Id: id}).
		Return(createWallets(), nil)

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)

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

	profileService := profile.NewMockProfileServiceClient(ctrl)
	profileService.EXPECT().
		GetWallets(ctx, &profileGen.UserID{Id: id}).
		Return(&profileGen.Wallets{}, errors.New("error"))

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)

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

	profileService := profile.NewMockProfileServiceClient(ctrl)
	profileService.EXPECT().
		CreateWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Empty{}, nil)

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)

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

	profileService := profile.NewMockProfileServiceClient(ctrl)
	profileService.EXPECT().
		CreateWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Empty{}, errors.New("error"))

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestSetTransaction_Success(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createSetTransactionSuccess(t, ctx)
	defer ctrl.Finish()

	desc, err := testAuth.SetTransaction(ctx, entity.Payment{
		UserId:   id,
		Base:     curr,
		Currency: base,
		Amount:   decimal.NewFromFloat(amount),
		Sell:     false,
	})

	require.NoError(t, err)
	require.Empty(t, desc)
}

func createSetTransactionSuccess(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewMockProfileServiceClient(ctrl)
	profileService.EXPECT().
		CheckWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Check{Existence: true}, nil)
	profileService.EXPECT().
		GetWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Wallet{Title: base, Value: 100 * amount}, nil)
	profileService.EXPECT().
		CheckWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: base}).
		Return(&profileGen.Check{Existence: true}, nil)
	profileService.EXPECT().
		UpdateWallet(ctx, &profileGen.ToSetWallet{Id: id, NewWallet: &profileGen.Wallet{Title: curr, Value: -amount * currValue}}).
		Return(&profileGen.Empty{}, nil)
	profileService.EXPECT().
		UpdateWallet(ctx, &profileGen.ToSetWallet{Id: id, NewWallet: &profileGen.Wallet{Title: base, Value: amount}}).
		Return(&profileGen.Empty{}, nil)
	profileService.EXPECT().
		AddToPaymentHistory(ctx, &profileGen.AddToHistory{
			Id:          id,
			Transaction: &profileGen.PaymentHistory{
				Currency: base,
				Base: curr,
				Value: currValue,
				Amount: amount,
				Sell: "false",
			},
		}).Return(&profileGen.Empty{}, nil)

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: curr}).Return(&currencyGen.Currency{Value: currValue}, nil)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: base}).Return(&currencyGen.Currency{Value: baseValue}, nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestSetTransaction_FailAddToHistory(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createSetTransactionFailAddToHistory(t, ctx)
	defer ctrl.Finish()

	desc, err := testAuth.SetTransaction(ctx, entity.Payment{
		UserId:   id,
		Base:     curr,
		Currency: base,
		Amount:   decimal.NewFromFloat(amount),
		Sell:     false,
	})

	require.Error(t, err)
	require.NotEmpty(t, desc)
}

func createSetTransactionFailAddToHistory(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewMockProfileServiceClient(ctrl)
	profileService.EXPECT().
		CheckWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Check{Existence: true}, nil)
	profileService.EXPECT().
		GetWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Wallet{Title: base, Value: 100 * amount}, nil)
	profileService.EXPECT().
		CheckWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: base}).
		Return(&profileGen.Check{Existence: true}, nil)
	profileService.EXPECT().
		UpdateWallet(ctx, &profileGen.ToSetWallet{Id: id, NewWallet: &profileGen.Wallet{Title: curr, Value: -amount * currValue}}).
		Return(&profileGen.Empty{}, nil)
	profileService.EXPECT().
		UpdateWallet(ctx, &profileGen.ToSetWallet{Id: id, NewWallet: &profileGen.Wallet{Title: base, Value: amount}}).
		Return(&profileGen.Empty{}, nil)
	profileService.EXPECT().
		AddToPaymentHistory(ctx, &profileGen.AddToHistory{
			Id:          id,
			Transaction: &profileGen.PaymentHistory{
				Currency: base,
				Base: curr,
				Value: currValue,
				Amount: amount,
				Sell: "false",
			},
		}).Return(&profileGen.Empty{}, errors.New("error"))

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: curr}).Return(&currencyGen.Currency{Value: currValue}, nil)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: base}).Return(&currencyGen.Currency{Value: baseValue}, nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestSetTransaction_FailUpdatePutWallet(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createSetTransactionFailUpdatePutWallet(t, ctx)
	defer ctrl.Finish()

	desc, err := testAuth.SetTransaction(ctx, entity.Payment{
		UserId:   id,
		Base:     curr,
		Currency: base,
		Amount:   decimal.NewFromFloat(amount),
		Sell:     false,
	})

	require.Error(t, err)
	require.NotEmpty(t, desc)
}

func createSetTransactionFailUpdatePutWallet(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewMockProfileServiceClient(ctrl)
	profileService.EXPECT().
		CheckWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Check{Existence: true}, nil)
	profileService.EXPECT().
		GetWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Wallet{Title: base, Value: 100 * amount}, nil)
	profileService.EXPECT().
		CheckWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: base}).
		Return(&profileGen.Check{Existence: true}, nil)
	profileService.EXPECT().
		UpdateWallet(ctx, &profileGen.ToSetWallet{Id: id, NewWallet: &profileGen.Wallet{Title: curr, Value: -amount * currValue}}).
		Return(&profileGen.Empty{}, nil)
	profileService.EXPECT().
		UpdateWallet(ctx, &profileGen.ToSetWallet{Id: id, NewWallet: &profileGen.Wallet{Title: base, Value: amount}}).
		Return(&profileGen.Empty{}, errors.New("error"))

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: curr}).Return(&currencyGen.Currency{Value: currValue}, nil)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: base}).Return(&currencyGen.Currency{Value: baseValue}, nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestSetTransaction_FailUpdateRemoveWallet(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createSetTransactionFailUpdateRemoveWallet(t, ctx)
	defer ctrl.Finish()

	desc, err := testAuth.SetTransaction(ctx, entity.Payment{
		UserId:   id,
		Base:     curr,
		Currency: base,
		Amount:   decimal.NewFromFloat(amount),
		Sell:     false,
	})

	require.Error(t, err)
	require.NotEmpty(t, desc)
}

func createSetTransactionFailUpdateRemoveWallet(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewMockProfileServiceClient(ctrl)
	profileService.EXPECT().
		CheckWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Check{Existence: true}, nil)
	profileService.EXPECT().
		GetWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Wallet{Title: base, Value: 100 * amount}, nil)
	profileService.EXPECT().
		CheckWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: base}).
		Return(&profileGen.Check{Existence: true}, nil)
	profileService.EXPECT().
		UpdateWallet(ctx, &profileGen.ToSetWallet{Id: id, NewWallet: &profileGen.Wallet{Title: curr, Value: -amount * currValue}}).
		Return(&profileGen.Empty{}, errors.New("error"))

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: curr}).Return(&currencyGen.Currency{Value: currValue}, nil)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: base}).Return(&currencyGen.Currency{Value: baseValue}, nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestSetTransaction_FailCheckWalletBuy(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createSetTransactionFailCheckWalletBuy(t, ctx)
	defer ctrl.Finish()

	desc, err := testAuth.SetTransaction(ctx, entity.Payment{
		UserId:   id,
		Base:     curr,
		Currency: base,
		Amount:   decimal.NewFromFloat(amount),
		Sell:     false,
	})

	require.Error(t, err)
	require.NotEmpty(t, desc)
}

func createSetTransactionFailCheckWalletBuy(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewMockProfileServiceClient(ctrl)
	profileService.EXPECT().
		CheckWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Check{Existence: true}, nil)
	profileService.EXPECT().
		GetWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Wallet{Title: base, Value: 100 * amount}, nil)
	profileService.EXPECT().
		CheckWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: base}).
		Return(&profileGen.Check{Existence: true}, errors.New("error"))

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: curr}).Return(&currencyGen.Currency{Value: currValue}, nil)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: base}).Return(&currencyGen.Currency{Value: baseValue}, nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestSetTransaction_SuccessCreateWallet(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createSetTransactionSuccessCreateWallet(t, ctx)
	defer ctrl.Finish()

	desc, err := testAuth.SetTransaction(ctx, entity.Payment{
		UserId:   id,
		Base:     curr,
		Currency: base,
		Amount:   decimal.NewFromFloat(amount),
		Sell:     false,
	})

	require.NoError(t, err)
	require.Empty(t, desc)
}

func createSetTransactionSuccessCreateWallet(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewMockProfileServiceClient(ctrl)
	profileService.EXPECT().
		CheckWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Check{Existence: true}, nil)
	profileService.EXPECT().
		GetWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Wallet{Title: base, Value: 100 * amount}, nil)
	profileService.EXPECT().
		CheckWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: base}).
		Return(&profileGen.Check{Existence: false}, nil)
	profileService.EXPECT().
		CreateWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: base}).
		Return(&profileGen.Empty{}, nil)
	profileService.EXPECT().
		UpdateWallet(ctx, &profileGen.ToSetWallet{Id: id, NewWallet: &profileGen.Wallet{Title: curr, Value: -amount * currValue}}).
		Return(&profileGen.Empty{}, nil)
	profileService.EXPECT().
		UpdateWallet(ctx, &profileGen.ToSetWallet{Id: id, NewWallet: &profileGen.Wallet{Title: base, Value: amount}}).
		Return(&profileGen.Empty{}, nil)
	profileService.EXPECT().
		AddToPaymentHistory(ctx, &profileGen.AddToHistory{
			Id:          id,
			Transaction: &profileGen.PaymentHistory{
				Currency: base,
				Base: curr,
				Value: currValue,
				Amount: amount,
				Sell: "false",
			},
		}).Return(&profileGen.Empty{}, nil)

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: curr}).Return(&currencyGen.Currency{Value: currValue}, nil)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: base}).Return(&currencyGen.Currency{Value: baseValue}, nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestSetTransaction_FailCreateWallet(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createSetTransactionFailCreateWallet(t, ctx)
	defer ctrl.Finish()

	desc, err := testAuth.SetTransaction(ctx, entity.Payment{
		UserId:   id,
		Base:     curr,
		Currency: base,
		Amount:   decimal.NewFromFloat(amount),
		Sell:     false,
	})

	require.Error(t, err)
	require.NotEmpty(t, desc)
}

func createSetTransactionFailCreateWallet(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewMockProfileServiceClient(ctrl)
	profileService.EXPECT().
		CheckWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Check{Existence: true}, nil)
	profileService.EXPECT().
		GetWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Wallet{Title: base, Value: 100 * amount}, nil)
	profileService.EXPECT().
		CheckWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: base}).
		Return(&profileGen.Check{Existence: false}, nil)
	profileService.EXPECT().
		CreateWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: base}).
		Return(&profileGen.Empty{}, errors.New("error"))

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: curr}).Return(&currencyGen.Currency{Value: currValue}, nil)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: base}).Return(&currencyGen.Currency{Value: baseValue}, nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestSetTransaction_FailGetWallet(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createSetTransactionFailGetWallet(t, ctx)
	defer ctrl.Finish()

	desc, err := testAuth.SetTransaction(ctx, entity.Payment{
		UserId:   id,
		Base:     curr,
		Currency: base,
		Amount:   decimal.NewFromFloat(amount),
		Sell:     false,
	})

	require.Error(t, err)
	require.NotEmpty(t, desc)
}

func createSetTransactionFailGetWallet(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewMockProfileServiceClient(ctrl)
	profileService.EXPECT().
		CheckWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Check{Existence: true}, nil)
	profileService.EXPECT().
		GetWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Wallet{}, errors.New("error"))

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: curr}).Return(&currencyGen.Currency{Value: currValue}, nil)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: base}).Return(&currencyGen.Currency{Value: baseValue}, nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestSetTransaction_FailNoPayment(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createSetTransactionFailNoPayment(t, ctx)
	defer ctrl.Finish()

	desc, err := testAuth.SetTransaction(ctx, entity.Payment{
		UserId:   id,
		Base:     curr,
		Currency: base,
		Amount:   decimal.NewFromFloat(amount),
		Sell:     false,
	})

	require.Error(t, err)
	require.NotEmpty(t, desc)
}

func createSetTransactionFailNoPayment(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewMockProfileServiceClient(ctrl)
	profileService.EXPECT().
		CheckWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Check{Existence: true}, nil)
	profileService.EXPECT().
		GetWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Wallet{Title: base, Value: amount}, nil)

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: curr}).Return(&currencyGen.Currency{Value: currValue}, nil)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: base}).Return(&currencyGen.Currency{Value: baseValue}, nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestSetTransaction_FailCheckWalletSell(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createSetTransactionFailCheckWalletSell(t, ctx)
	defer ctrl.Finish()

	desc, err := testAuth.SetTransaction(ctx, entity.Payment{
		UserId:   id,
		Base:     curr,
		Currency: base,
		Amount:   decimal.NewFromFloat(amount),
		Sell:     false,
	})

	require.Error(t, err)
	require.NotEmpty(t, desc)
}

func createSetTransactionFailCheckWalletSell(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewMockProfileServiceClient(ctrl)
	profileService.EXPECT().
		CheckWallet(ctx, &profileGen.ConcreteWallet{Id: id, Title: curr}).
		Return(&profileGen.Check{Existence: true}, errors.New("error"))

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: curr}).Return(&currencyGen.Currency{Value: currValue}, nil)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: base}).Return(&currencyGen.Currency{Value: baseValue}, nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestSetTransaction_FailGetLastRateBase(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createSetTransactionFailGetLastRateBase(t, ctx)
	defer ctrl.Finish()

	desc, err := testAuth.SetTransaction(ctx, entity.Payment{
		UserId:   id,
		Base:     curr,
		Currency: base,
		Amount:   decimal.NewFromFloat(amount),
		Sell:     false,
	})

	require.Error(t, err)
	require.NotEmpty(t, desc)
}

func createSetTransactionFailGetLastRateBase(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewMockProfileServiceClient(ctrl)

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: curr}).Return(&currencyGen.Currency{Value: currValue}, nil)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: base}).Return(&currencyGen.Currency{}, errors.New("error"))

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestSetTransaction_FailGetLastRateCurr(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createSetTransactionFailGetLastRateCurr(t, ctx)
	defer ctrl.Finish()

	desc, err := testAuth.SetTransaction(ctx, entity.Payment{
		UserId:   id,
		Base:     curr,
		Currency: base,
		Amount:   decimal.NewFromFloat(amount),
		Sell:     false,
	})

	require.Error(t, err)
	require.NotEmpty(t, desc)
}

func createSetTransactionFailGetLastRateCurr(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewMockProfileServiceClient(ctrl)

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	currencyService.EXPECT().
		GetLastRate(ctx, &currencyGen.CurrencyTitle{Title: curr}).Return(&currencyGen.Currency{}, errors.New("error"))

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestGetAllIncomePerDay_Success(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createGetAllIncomePerDaySuccess(t, ctx)
	defer ctrl.Finish()

	desc, walletStates, err := testAuth.GetAllIncomePerDay(ctx, id)

	require.NoError(t, err)
	require.Empty(t, desc)
	require.NotEmpty(t, walletStates)
}

func createGetAllIncomePerDaySuccess(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewMockProfileServiceClient(ctrl)

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	profileService.EXPECT().
		GetAllIncomePerDay(ctx, &profileGen.UserID{Id: id}).
		Return(createWalletStates(), nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewPaymentApp(profileService, currencyService, securityService), ctrl
}

func TestGetAllIncomePerDay_Fail(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createGetAllIncomePerDayFail(t, ctx)
	defer ctrl.Finish()

	desc, walletStates, err := testAuth.GetAllIncomePerDay(ctx, id)

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Equal(t, "GetAllIncomePerDay", desc.Action)
	require.Equal(t, http.StatusInternalServerError, desc.Status)
	require.Empty(t, walletStates)
}

func createGetAllIncomePerDayFail(t *testing.T, ctx context.Context) (*application.PaymentApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewMockProfileServiceClient(ctrl)

	currencyService := currency.NewMockCurrencyServiceClient(ctrl)
	profileService.EXPECT().
		GetAllIncomePerDay(ctx, &profileGen.UserID{Id: id}).
		Return(&profileGen.WalletStates{}, errors.New("fail"))

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

func createWalletStates() *profileGen.WalletStates {
	return &profileGen.WalletStates{
		States: []*profileGen.WalletState{
			{Value: currValue, UpdatedAt: &timestamppb.Timestamp{Seconds: 100, Nanos: 100}},
	}}
}
