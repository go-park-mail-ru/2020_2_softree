package repository

import "server/src/canal/domain/entity"

type DayCurrencyRepository interface {
	dayCurrencySaver
	dayCurrencyReceiver
}

type dayCurrencySaver interface {
	SaveCurrency(financial FinancialRepository) error
}

type dayCurrencyReceiver interface {
	GetInitialCurrency() ([]entity.Currency, error)
}
