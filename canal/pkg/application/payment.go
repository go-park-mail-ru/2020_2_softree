package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"net/http"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/domain/repository"
	currency "server/currency/pkg/currency/gen"
	profile "server/profile/pkg/profile/gen"
	"time"
)

type PaymentApp struct {
	profile   profile.ProfileServiceClient
	currency  currency.CurrencyServiceClient
	sanitizer bluemonday.Policy
	security  repository.Utils
}

func NewPaymentApp(profile profile.ProfileServiceClient, rates currency.CurrencyServiceClient, security repository.Utils) *PaymentApp {
	return &PaymentApp{profile, rates, *bluemonday.UGCPolicy(), security}
}

func (pmt *PaymentApp) ReceiveTransactions(ctx context.Context, in entity.Income) (entity.Description, entity.Payments, error) {
	history, err := pmt.profile.GetAllPaymentHistory(ctx, &profile.IncomeParameters{Id: in.Id, Period: in.Period})
	if err != nil {
		return createErrorDescription("ReceiveTransactions", "GetAllPaymentHistory", http.StatusInternalServerError),
			entity.Payments{}, err
	}

	return entity.Description{}, entity.ConvertToPayment(history), err
}

func (pmt *PaymentApp) ReceiveWallets(ctx context.Context, id int64) (entity.Description, entity.Wallets, error) {
	wallets, err := pmt.profile.GetWallets(ctx, &profile.UserID{Id: id})
	if err != nil {
		return createErrorDescription("ReceiveWallets", "GetWallets", http.StatusInternalServerError),
			entity.Wallets{}, err
	}

	return entity.Description{}, entity.ConvertToWallets(wallets), nil
}

func (pmt *PaymentApp) SetWallet(ctx context.Context, wallet entity.Wallet) (entity.Description, error) {
	if _, err := pmt.profile.CreateWallet(ctx, &profile.ConcreteWallet{Id: wallet.UserId, Title: wallet.Title}); err != nil {
		return createErrorDescription("SetWallet", "CreateWallet", http.StatusInternalServerError), err
	}

	return entity.Description{}, nil
}

func (pmt *PaymentApp) SetTransaction(ctx context.Context, payment entity.Payment) (entity.Description, error) {
	var div decimal.Decimal
	var desc entity.Description
	var err error
	transaction := payment.ConvertToGRPC()
	if desc, div, err = pmt.getCurrencyDiv(ctx, transaction); err != nil {
		return desc, err
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

	var exist bool
	if exist, desc, err = pmt.checkWalletSell(ctx, &profile.ConcreteWallet{Id: payment.UserId, Title: titleToCheckPayment}); !exist {
		return desc, err
	}

	if desc, err = pmt.getPay(ctx, &profile.ConcreteWallet{Id: payment.UserId, Title: titleToCheckPayment}, checkingPayment); err != nil {
		if desc.Status == notEnoughPayment {
			errs := entity.ErrorJSON{}
			errs.NonFieldError = append(errs.NonFieldError, "not enough payment")
			desc.ErrorJSON = errs
			desc.Status = http.StatusBadRequest
		}
		return desc, err
	}

	if exist, desc, err = pmt.checkWalletBuy(ctx, &profile.ConcreteWallet{Id: payment.UserId, Title: putTitle}); !exist {
		return desc, err
	}

	toSetWallet := profile.ToSetWallet{Id: payment.UserId, NewWallet: &profile.Wallet{Title: removedTitle, Value: removedMoney}}
	if _, err = pmt.profile.UpdateWallet(ctx, &toSetWallet); err != nil {
		return createErrorDescription("SetTransactions", "UpdateWallet "+toSetWallet.NewWallet.Title, http.StatusInternalServerError), err
	}

	toSetWallet = profile.ToSetWallet{Id: payment.UserId, NewWallet: &profile.Wallet{Title: putTitle, Value: putMoney}}
	if _, err := pmt.profile.UpdateWallet(ctx, &toSetWallet); err != nil {
		return createErrorDescription("SetTransactions", "UpdateWallet "+toSetWallet.NewWallet.Title, http.StatusInternalServerError), err
	}

	transaction.Value, _ = div.Float64()
	if _, err = pmt.profile.AddToPaymentHistory(ctx, &profile.AddToHistory{Id: payment.UserId, Transaction: transaction}); err != nil {
		return createErrorDescription("SetTransactions", "AddToPaymentHistory", http.StatusInternalServerError), err
	}

	return entity.Description{}, nil
}

func (pmt *PaymentApp) checkWalletSell(ctx context.Context, wallet *profile.ConcreteWallet) (bool, entity.Description, error) {
	var exist *profile.Check
	var err error
	if exist, err = pmt.profile.CheckWallet(ctx, wallet); err != nil {
		return false, createErrorDescription("checkWalletSell", "CheckWallet", http.StatusInternalServerError), err
	}

	if !exist.Existence {
		return false, entity.Description{Status: http.StatusBadRequest}, errors.New("existence")
	}

	return true, entity.Description{}, nil
}

func (pmt *PaymentApp) checkWalletBuy(ctx context.Context, wallet *profile.ConcreteWallet) (bool, entity.Description, error) {
	var exist *profile.Check
	var err error
	if exist, err = pmt.profile.CheckWallet(ctx, wallet); err != nil {
		return false, createErrorDescription("checkWalletBuy", "CheckWallet", http.StatusInternalServerError), err
	}

	if !exist.Existence {
		if _, err = pmt.profile.CreateWallet(ctx, wallet); err != nil {
			return false, createErrorDescription("checkWalletBuy", "CreateWallet", http.StatusInternalServerError), err
		}
	}

	return true, entity.Description{}, nil
}

func (pmt *PaymentApp) getCurrencyDiv(ctx context.Context, transaction *profile.PaymentHistory) (entity.Description, decimal.Decimal, error) {
	var currencyBase *currency.Currency
	var err error
	if currencyBase, err = pmt.currency.GetLastRate(ctx, &currency.CurrencyTitle{Title: transaction.Base}); err != nil {
		return createErrorDescription("getCurrencyDiv", "GetLastRate: Base", http.StatusInternalServerError),
			decimal.Decimal{}, err
	}

	var currencyCurr *currency.Currency
	if currencyCurr, err = pmt.currency.GetLastRate(ctx, &currency.CurrencyTitle{Title: transaction.Currency}); err != nil {
		return createErrorDescription("getCurrencyDiv", "GetLastRate: Currency", http.StatusInternalServerError),
			decimal.Decimal{}, err
	}

	div := decimal.NewFromFloat(currencyBase.Value).Div(decimal.NewFromFloat(currencyCurr.Value))
	return entity.Description{}, div, nil
}

const notEnoughPayment = 1

func (pmt *PaymentApp) getPay(ctx context.Context, userWallet *profile.ConcreteWallet, needToPay decimal.Decimal) (entity.Description, error) {
	var wallet *profile.Wallet
	var err error
	if wallet, err = pmt.profile.GetWallet(ctx, userWallet); err != nil {
		return createErrorDescription("getPay", "GetWallet", http.StatusInternalServerError), err
	}

	if needToPay.GreaterThan(decimal.NewFromFloat(wallet.Value)) {
		return entity.Description{Status: notEnoughPayment}, fmt.Errorf("%d notEnoughPayment", userWallet.Id)
	}

	return entity.Description{}, nil
}

func (pmt *PaymentApp) GetIncome(ctx context.Context, in entity.Income) (entity.Description, decimal.Decimal, error) {
	incomeParameters := in.ConvertToGRPC()
	result, err := pmt.profile.GetIncome(ctx, incomeParameters)
	if err != nil {
		return createErrorDescription("GetIncome", "GetIncome", http.StatusInternalServerError), decimal.Decimal{}, err
	}

	walletUSDCash, err := pmt.transformActualUserWallets(ctx, incomeParameters.Id)
	if err != nil {
		return createErrorDescription("GetIncome", "transformActualUserWallets", http.StatusInternalServerError),
			decimal.Decimal{}, err
	}

	return entity.Description{}, walletUSDCash.Sub(decimal.NewFromFloat(result.Change)).Round(3), nil
}

func (pmt *PaymentApp) GetAllIncomePerDay(ctx context.Context, in entity.Income) (entity.Description, entity.WalletStates, error) {
	out, err := pmt.profile.GetAllIncomePerDay(ctx, &profile.IncomeParameters{Id: in.Id, Period: in.Period})
	if err != nil {
		return createErrorDescription("GetAllIncomePerDay", "GetAllIncomePerDay", http.StatusInternalServerError),
			entity.WalletStates{}, err
	}

	return entity.Description{}, entity.ConvertToWalletStates(out), nil
}

func (pmt *PaymentApp) transformActualUserWallets(ctx context.Context, id int64) (decimal.Decimal, error) {
	wallets, err := pmt.profile.GetWallets(ctx, &profile.UserID{Id: id})
	if err != nil {
		return decimal.Decimal{}, err
	}

	var cash decimal.Decimal
	for _, wallet := range wallets.Wallets {
		curr, err := pmt.currency.GetLastRate(ctx, &currency.CurrencyTitle{Title: wallet.Title})
		if err != nil {
			return decimal.Decimal{}, err
		}
		cash = cash.Add(decimal.NewFromFloat(wallet.Value / curr.Value))
	}

	return cash, nil
}

// every day task
func (pmt *PaymentApp) WritePortfolios() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userNum, err := pmt.profile.GetUsers(ctx, &profile.Empty{})
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "WritePortfolio", "action": "GetUsers"}).Error(err)
		return
	}

	for i := int64(0); i < userNum.Num; i++ {
		portfolioValue, err := pmt.transformActualUserWallets(ctx, i)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"function": "WritePortfolio",
				"action":   "transformActualUserWallets",
				"user_id":  i,
			}).Debug(err)

			return
		}

		value, _ := portfolioValue.Float64()
		_, err = pmt.profile.PutPortfolio(ctx, &profile.PortfolioValue{Id: i, Value: value})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"function": "WritePortfolio",
				"action":   "PutPortfolio",
				"user_id":  i,
				"value":    value,
			}).Debug(err)

			return
		}
	}
}
