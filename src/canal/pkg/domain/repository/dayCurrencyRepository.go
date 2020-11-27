package repository

import (
	"server/src/canal/pkg/domain/entity"
	"server/src/currency/pkg/domain"
)

type DayCurrencyRepository interface {
	dayCurrencySaver
	dayCurrencyReceiver
}

type dayCurrencySaver interface {
	SaveCurrency(financial domain.FinancialRepository) error
}

type dayCurrencyReceiver interface {
	GetInitialCurrency() ([]entity.Currency, error)
}
