package domain

type FinancialRepository interface {
	base
	currency
}

type currency interface {
	GetQuote() map[string]interface{}
}

type base interface {
	GetBase() string
}
