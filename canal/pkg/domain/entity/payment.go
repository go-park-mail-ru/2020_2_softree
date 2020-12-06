package entity

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/shopspring/decimal"
	profile "server/profile/pkg/profile/gen"
	"time"
)

type Payment struct {
	Currency  string          `json:"currency"`
	Base      string          `json:"base"`
	Amount    decimal.Decimal `json:"amount"`
	Value     decimal.Decimal `json:"value"`
	Sell      bool            `json:"sell"`
	UpdatedUp time.Time       `json:"updated_up"`
}

func ConvertToPayment(history *profile.AllHistory) []Payment {
	payments := make([]Payment, 0, len(history.History))
	for _, pay := range history.History {
		updated, _ := ptypes.Timestamp(pay.UpdatedAt)
		payments = append(payments, Payment{
			Currency:  pay.Currency,
			Base:      pay.Base,
			Amount:    decimal.NewFromFloat(pay.Amount),
			Value:     decimal.NewFromFloat(pay.Value),
			Sell:      !(pay.Sell == "false"),
			UpdatedUp: updated,
		})
	}

	return payments
}
