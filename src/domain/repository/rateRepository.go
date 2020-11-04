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
	SaveRate(entity.Rate) (entity.Rate, error)
}

type rateUpdater interface {
	UpdateRate(uint64, entity.Rate) (entity.Rate, error)
}

type rateEraser interface {
	DeleteRate(uint64) error
}

type ratesReceiver interface {
	GetRates() ([]entity.Rate, error)
}

type rateReceiver interface {
	GetRate(uint64) (entity.Rate, error)
}
