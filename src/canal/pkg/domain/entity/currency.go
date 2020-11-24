package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type Currency struct {
	Title     string          `json:"title"`
	Value     decimal.Decimal `json:"value,omitempty"`
	UpdatedAt time.Time       `json:"updated_at,omitempty"`
	Base      string          `json:"base,omitempty"`
}
