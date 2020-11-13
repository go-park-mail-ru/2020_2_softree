package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type PaymentHistory struct {
	From     string          `json:"from"`
	To       string          `json:"to"`
	Amount   decimal.Decimal `json:"amount"`
	Datetime time.Time       `json:"datetime"`
	Value    decimal.Decimal `json:"value"`
}

type Interval struct {
	From  time.Time `json:"from"`
	Where time.Time `json:"where"`
}
