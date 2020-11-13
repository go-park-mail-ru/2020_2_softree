package entity

import "github.com/shopspring/decimal"

type Wallet struct {
	Title string          `json:"title"`
	Value decimal.Decimal `json:"value"`
}
