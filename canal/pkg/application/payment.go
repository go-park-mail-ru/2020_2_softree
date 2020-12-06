package application

import (
	"context"
	"errors"
	"github.com/microcosm-cc/bluemonday"
	"github.com/shopspring/decimal"
	"net/http"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/domain/repository"
	currency "server/currency/pkg/currency/gen"
	profile "server/profile/pkg/profile/gen"
)

type PaymentApp struct {
	profile   profile.ProfileServiceClient
	rates     currency.CurrencyServiceClient
	sanitizer bluemonday.Policy
	security  repository.Utils
}

func NewPaymentApp(profile profile.ProfileServiceClient, security repository.Utils) *PaymentApp {
	return &PaymentApp{profile: profile, security: security, sanitizer: *bluemonday.UGCPolicy()}
}

func (pmt *PaymentApp) ReceiveTransactions(ctx context.Context, id int64) (entity.Description, entity.Payments) {
	history, err := pmt.profile.GetAllPaymentHistory(ctx, &profile.UserID{Id: id})
	if err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "ReceiveTransactions",
			Action:   "GetAllPaymentHistory",
			Err:      err,
		}, entity.Payments{}
	}

	return entity.Description{}, entity.ConvertToPayment(history)
}

func (pmt *PaymentApp) SetTransaction(ctx context.Context, payment entity.Payment) entity.Description {
	var div decimal.Decimal
	var desc entity.Description
	transaction := payment.ConvertToGRPC()
	if desc, div = pmt.getCurrencyDiv(ctx, transaction); desc.Err != nil {
		return desc
	}

	divMulAmount := div.Mul(decimal.NewFromFloat(transaction.Amount))

	titleToCheckPayment := transaction.Currency
	checkingPayment := decimal.NewFromFloat(transaction.Amount)

	removedMoney := -transaction.Amount
	removedTitle := transaction.Currency

	putMoney, _ := divMulAmount.Float64()
	putTitle := transaction.Base
	if !payment.Sell {
		titleToCheckPayment = transaction.Base
		checkingPayment = divMulAmount

		removedMoney, _ = checkingPayment.Float64()
		removedMoney *= -1
		removedTitle = transaction.Base

		putMoney = transaction.Amount
		putTitle = transaction.Currency
	}

	if exist, desc := pmt.checkWalletSell(ctx, &profile.ConcreteWallet{Id: payment.UserId, Title: titleToCheckPayment}); !exist {
		return desc
	}

	if desc := pmt.getPay(ctx, &profile.ConcreteWallet{Id: payment.UserId, Title: titleToCheckPayment}, checkingPayment); desc.Err != nil {
		if desc.Status == notEnoughPayment {
			errs := entity.ErrorJSON{}
			errs.NonFieldError = append(errs.NonFieldError, "not enough payment")
			desc.ErrorJSON = errs
		}
		return desc
	}

	if exist, desc := pmt.checkWalletBuy(ctx, &profile.ConcreteWallet{Id: payment.UserId, Title: putTitle}); !exist {
		return desc
	}

	toSetWallet := profile.ToSetWallet{Id: payment.UserId, NewWallet: &profile.Wallet{Title: removedTitle, Value: removedMoney}}
	if _, err := pmt.profile.UpdateWallet(ctx, &toSetWallet); err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "SetTransactions",
			Action:   "UpdateWallet " + toSetWallet.NewWallet.Title,
			Err:      err,
		}
	}

	toSetWallet = profile.ToSetWallet{Id: payment.UserId, NewWallet: &profile.Wallet{Title: putTitle, Value: putMoney}}
	if _, err := pmt.profile.UpdateWallet(ctx, &toSetWallet); err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "SetTransactions",
			Action:   "UpdateWallet" + toSetWallet.NewWallet.Title,
			Err:      err,
		}
	}

	transaction.Value, _ = div.Float64()
	if _, err := pmt.profile.AddToPaymentHistory(ctx, &profile.AddToHistory{Id: payment.UserId, Transaction: transaction}); err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "SetTransactions",
			Action:   "AddToPaymentHistory",
			Err:      err,
		}
	}

	return entity.Description{}
}

func (pmt *PaymentApp) checkWalletSell(ctx context.Context, wallet *profile.ConcreteWallet) (bool, entity.Description) {
	var exist *profile.Check
	var err error
	if exist, err = pmt.profile.CheckWallet(ctx, wallet); err != nil {
		return false, entity.Description{
			Status:   http.StatusInternalServerError,
			Err:      err,
			Function: "checkWalletSell",
			Action:   "CheckWallet",
		}
	}

	if !exist.Existence {
		return false, entity.Description{
			Status: http.StatusBadRequest,
			Err:    errors.New("existence"),
		}
	}

	return true, entity.Description{}
}

func (pmt *PaymentApp) checkWalletBuy(ctx context.Context, wallet *profile.ConcreteWallet) (bool, entity.Description) {
	var exist *profile.Check
	var err error
	if exist, err = pmt.profile.CheckWallet(ctx, wallet); err != nil {
		return false, entity.Description{
			Function: "checkWalletBuy",
			Status:   http.StatusInternalServerError,
			Action:   "CheckWallet",
			Err:      err,
		}
	}

	if !exist.Existence {
		if _, err = pmt.profile.CreateWallet(ctx, wallet); err != nil {
			return false, entity.Description{
				Function: "checkWalletBuy",
				Status:   http.StatusInternalServerError,
				Action:   "CreateWallet",
				Err:      err,
			}
		}
	}

	return true, entity.Description{}
}

func (pmt *PaymentApp) getCurrencyDiv(ctx context.Context, transaction *profile.PaymentHistory) (entity.Description, decimal.Decimal) {
	var currencyBase *currency.Currency
	var err error
	if currencyBase, err = pmt.rates.GetLastRate(ctx, &currency.CurrencyTitle{Title: transaction.Base}); err != nil {
		return entity.Description{
			Err:      err,
			Status:   http.StatusInternalServerError,
			Function: "getCurrencyDiv",
			Action:   "GetLastRate",
		}, decimal.Decimal{}
	}

	var currencyCurr *currency.Currency
	if currencyCurr, err = pmt.rates.GetLastRate(ctx, &currency.CurrencyTitle{Title: transaction.Currency}); err != nil {
		return entity.Description{
			Err:      err,
			Status:   http.StatusInternalServerError,
			Function: "getCurrencyDiv",
			Action:   "GetLastRate",
		}, decimal.Decimal{}
	}

	div := decimal.NewFromFloat(currencyBase.Value).Div(decimal.NewFromFloat(currencyCurr.Value))
	return entity.Description{}, div
}

const notEnoughPayment = 1

func (pmt *PaymentApp) getPay(ctx context.Context, userWallet *profile.ConcreteWallet, needToPay decimal.Decimal) entity.Description {
	var wallet *profile.Wallet
	var err error
	if wallet, err = pmt.profile.GetWallet(ctx, userWallet); err != nil {
		return entity.Description{
			Function: "getPay",
			Status:   http.StatusInternalServerError,
			Action:   "GetWallet",
			Err:      err,
		}
	}

	if needToPay.GreaterThan(decimal.NewFromFloat(wallet.Value)) {
		return entity.Description{
			Status: notEnoughPayment,
			Err:    nil,
		}
	}

	return entity.Description{}
}
