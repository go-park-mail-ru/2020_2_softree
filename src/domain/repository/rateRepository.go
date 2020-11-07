package repository

import "server/src/domain/entity"

type RateRepository interface {
	ratesSaver
	rateUpdater
	rateEraser
	ratesReceiver
	rateReceiver
}

type ratesSaver interface {
	SaveRates(financial FinancialRepository) error
}

type rateUpdater interface {
	UpdateRate(uint64, entity.Currency) (entity.Currency, error)
}

type rateEraser interface {
	DeleteRate(uint64) error
}

type ratesReceiver interface {
	GetRates() ([]entity.Currency, error)
}

type rateReceiver interface {
	GetRate(uint64) (entity.Currency, error)
}
