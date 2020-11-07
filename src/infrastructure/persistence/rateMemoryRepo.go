package persistence

import (
	"errors"
	"server/src/domain/entity"
	"server/src/domain/repository"
)

type RateMemoryRepo struct {
	rates     []entity.Rate
}

func NewRateRepository() *RateMemoryRepo {
	rates := make([]entity.Rate, 1)
	return &RateMemoryRepo{rates: rates}
}

func (rr *RateMemoryRepo) SaveRates(financial repository.FinancialRepository) error {
	for name, quote := range financial.GetQuote() {
		var rate entity.Rate

		rate.ID = uint64(len(rr.rates) + 1)
		rate.Base = financial.GetBase()
		rate.Currency = name
		rate.Value = quote.(float64)

		rr.rates = append(rr.rates, rate)
	}

	return nil
}

func (rr *RateMemoryRepo) UpdateRate(id uint64, data entity.Rate) (rate entity.Rate, err error) {
	return
}

func (rr *RateMemoryRepo) DeleteRate(id uint64) error {
	for i, rate := range rr.rates {
		if rate.ID == id {
			rr.rates = append(rr.rates[:i], rr.rates[i+1:]...)
		}
	}

	return nil
}

func (rr *RateMemoryRepo) GetRates() ([]entity.Rate, error) {
	return rr.rates, nil
}

func (rr *RateMemoryRepo) GetRate(id uint64) (rate entity.Rate, err error) {
	for _, rate = range rr.rates {
		if rate.ID == id {
			return rate, nil
		}
	}

	return entity.Rate{}, errors.New("no rate")
}
