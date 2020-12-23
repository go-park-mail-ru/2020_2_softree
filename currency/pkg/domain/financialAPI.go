package domain

type FinancialAPI interface {
	GetCurrencies() (FinancialRepository, error)
}
