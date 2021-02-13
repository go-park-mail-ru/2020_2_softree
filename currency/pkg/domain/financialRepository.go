package domain

type FinancialRepository interface {
	base
	currency
}

type currency interface {
	GetQuote() map[string]float64
}

type base interface {
	GetBase() string
}
