package application

import (
	"server/src/domain/entity"
	"server/src/domain/repository"
)

type RateApp struct {
	rr  repository.RateRepository
	dcr repository.DayCurrencyRepository
}

func NewRateApp(repo repository.RateRepository, dcr repository.DayCurrencyRepository) *RateApp {
	return &RateApp{rr: repo, dcr: dcr}
}

func (ra *RateApp) SaveCurrency(financial repository.FinancialRepository) error {
	return ra.dcr.SaveCurrency(financial)
}

func (ra *RateApp) GetInitialCurrency() ([]entity.Currency, error) {
	return ra.dcr.GetInitialCurrency()
}

func (ra *RateApp) SaveRates(financial repository.FinancialRepository) error {
	return ra.rr.SaveRates(financial)
}

func (ra *RateApp) UpdateRate(id uint64, rate entity.Currency) (entity.Currency, error) {
	return ra.rr.UpdateRate(id, rate)
}

func (ra *RateApp) DeleteRate(id uint64) error {
	return ra.rr.DeleteRate(id)
}

func (ra *RateApp) GetRates() ([]entity.Currency, error) {
	return ra.rr.GetRates()
}

func (ra *RateApp) GetRate(id uint64) (entity.Currency, error) {
	return ra.rr.GetRate(id)
}
