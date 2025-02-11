package repository

import (
	"context"
	"github.com/shopspring/decimal"
	"server/canal/pkg/domain/entity"
)

type PaymentLogic interface {
	ReceiveTransactions(ctx context.Context, in entity.Income) (entity.Description, entity.Payments, error)
	ReceiveWallets(ctx context.Context, id int64) (entity.Description, entity.Wallets, error)
	SetTransaction(ctx context.Context, payment entity.Payment) (entity.Description, error)
	SetWallet(ctx context.Context, wallet entity.Wallet) (entity.Description, error)
	GetIncome(ctx context.Context, in entity.Income) (entity.Description, decimal.Decimal, error)
	GetAllIncomePerDay(ctx context.Context, in entity.Income) (entity.Description, entity.WalletStates, error)
	WritePortfolios()
}
