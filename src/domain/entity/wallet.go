package entity

import "github.com/shopspring/decimal"

type Wallet struct {
	USD decimal.Decimal `json:"USD,omitempty"`
	RUB decimal.Decimal `json:"RUB,omitempty"`
	EUR decimal.Decimal `json:"EUR,omitempty"`
	GBP decimal.Decimal `json:"GBP,omitempty"`
}
