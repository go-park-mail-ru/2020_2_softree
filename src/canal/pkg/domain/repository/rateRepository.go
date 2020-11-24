package repository

import "server/src/canal/pkg/domain/entity"

type RateRepository interface {
	ratesSaver
	rateUpdater
	rateEraser
	ratesReceiver
	rateReceiver
	rateReceiverLast
}

type ratesSaver interface {
	SaveRates(financial FinancialRepository) error
}

type rateUpdater interface {
	UpdateRate(int64, entity.Currency) (entity.Currency, error)
}

type rateEraser interface {
	DeleteRate(uint64) error
}

type ratesReceiver interface {
	GetRates() ([]entity.Currency, error)
}

type rateReceiver interface {
	GetRate(string) ([]entity.Currency, error)
}

type rateReceiverLast interface {
	GetLastRate(string) (entity.Currency, error)
}
