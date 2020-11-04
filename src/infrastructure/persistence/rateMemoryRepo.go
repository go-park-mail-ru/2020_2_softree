package persistence

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"server/src/domain/entity"
)

type RateMemoryRepo struct {
	rates []entity.Rate
}

func NewRateRepository() *RateMemoryRepo {
	rates := make([]entity.Rate, 1)
	return &RateMemoryRepo{rates: rates}
}

func (rr *RateMemoryRepo) SaveRate(rate entity.Rate) (entity.Rate, error) {
	rate.ID = uint64(len(rr.rates) + 1)

	rr.rates = append(rr.rates, rate)
	return rate, nil
}

func (rr *RateMemoryRepo) UpdateRate(id uint64, data entity.Rate) (rate entity.Rate, err error) {
	var i int
	for i, rate = range rr.rates {
		if rate.ID == id {
			break
		}
	}

	if !govalidator.IsNull(data.Value) {
		rr.rates[i].Value = data.Value
		rate.Value = data.Value
	}

	return
}

func (rr *RateMemoryRepo) DeleteRate(id uint64) error {
	for i, rate := range rr.rates {
		if rate.ID == id {
			rr.rates = append(rr.rates[:i], rr.rates[i + 1:]...)
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
