package application

import (
	"server/src/domain/entity"
	"server/src/domain/repository"
)

type RateApp struct {
	rateRepository        repository.RateRepository
	dayCurrencyRepository repository.DayCurrencyRepository
}

func NewRateApp(repo repository.RateRepository, dcr repository.DayCurrencyRepository) *RateApp {
	return &RateApp{rateRepository: repo, dayCurrencyRepository: dcr}
}

func (ra *RateApp) SaveCurrency(financial repository.FinancialRepository) error {
	return ra.dayCurrencyRepository.SaveCurrency(financial)
}

func (ra *RateApp) GetInitialCurrency() ([]entity.Currency, error) {
	return ra.dayCurrencyRepository.GetInitialCurrency()
}

func (ra *RateApp) SaveRates(financial repository.FinancialRepository) error {
	return ra.rateRepository.SaveRates(financial)
}

func (ra *RateApp) UpdateRate(id uint64, rate entity.Currency) (entity.Currency, error) {
	return ra.rateRepository.UpdateRate(id, rate)
}

func (ra *RateApp) DeleteRate(id uint64) error {
	return ra.rateRepository.DeleteRate(id)
}

func (ra *RateApp) GetRates() ([]entity.Currency, error) {
	return ra.rateRepository.GetRates()
}

func (ra *RateApp) GetRate(title string) ([]entity.Currency, error) {
	return ra.rateRepository.GetRate(title)
}
