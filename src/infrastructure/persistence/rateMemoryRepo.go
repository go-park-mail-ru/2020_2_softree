package persistence

import (
	"errors"
	"server/src/domain/entity"
	"server/src/domain/repository"
)

type RateMemoryRepo struct {
	rates     []entity.Currency
}

func NewRateRepository() *RateMemoryRepo {
	rates := make([]entity.Currency, 1)
	return &RateMemoryRepo{rates: rates}
}

func (rr *RateMemoryRepo) SaveRates(financial repository.FinancialRepository) error {
	for name, quote := range financial.GetQuote() {
		var rate entity.Currency

		rate.ID = uint64(len(rr.rates) + 1)
		rate.Base = financial.GetBase()
		rate.Title = name
		rate.Value = quote.(float64)

		rr.rates = append(rr.rates, rate)
	}

	return nil
}

func (rr *RateMemoryRepo) UpdateRate(id uint64, data entity.Currency) (rate entity.Currency, err error) {
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

func (rr *RateMemoryRepo) GetRates() ([]entity.Currency, error) {
	return rr.rates, nil
}

func (rr *RateMemoryRepo) GetRate(id uint64) (rate entity.Currency, err error) {
	for _, rate = range rr.rates {
		if rate.ID == id {
			return rate, nil
		}
	}

	return entity.Currency{}, errors.New("no rate")
}
