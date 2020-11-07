package repository

import "server/src/domain/entity"

type DayCurrencyRepository interface {
	dayCurrencySaver
	dayCurrencyReceiver
}

type dayCurrencySaver interface {
	SaveCurrency([]entity.Currency) error
}

type dayCurrencyReceiver interface {
	GetInitialCurrency() ([]entity.Currency, error)
}
