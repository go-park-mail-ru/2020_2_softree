package application

import (
	"server/src/domain/entity"
	"server/src/domain/repository"
)

type RateApp struct {
	rr repository.RateRepository
}

func NewRateApp(repo repository.RateRepository) *RateApp {
	return &RateApp{rr: repo}
}

func (ra *RateApp) SaveRate(rate entity.Rate) (entity.Rate, error) {
	return ra.rr.SaveRate(rate)
}

func (ra *RateApp) UpdateRate(id uint64, rate entity.Rate) (entity.Rate, error) {
	return ra.rr.UpdateRate(id, rate)
}

func (ra *RateApp) DeleteRate(id uint64) error {
	return ra.rr.DeleteRate(id)
}

func (ra *RateApp) GetRates() ([]entity.Rate, error) {
	return ra.rr.GetRates()
}

func (ra *RateApp) GetRate(id uint64) (entity.Rate, error) {
	return ra.rr.GetRate(id)
}