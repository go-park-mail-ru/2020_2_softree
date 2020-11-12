package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type PaymentHistory struct {
	Base          string          `json:"base"`
	Title         string          `json:"title"`
	Value         decimal.Decimal `json:"value"`
	CurrencyValue decimal.Decimal `json:"currency_value"`
	Commission    decimal.Decimal `json:"commission"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type Interval struct {
	From  time.Time `json:"from"`
	Where time.Time `json:"where"`
}
