package application

import (
	"server/src/canal/domain/entity"
	"server/src/canal/domain/repository"
)

type RateApp struct {
	rateRepository        repository.RateRepository
	dayCurrencyRepository repository.DayCurrencyRepository
}

func NewRateApp(repo repository.RateRepository, dcr repository.DayCurrencyRepository) *RateApp {
	return &RateApp{rateRepository: repo, dayCurrencyRepository: dcr}
}

func (rateApp *RateApp) SaveCurrency(financial repository.FinancialRepository) error {
	return rateApp.dayCurrencyRepository.SaveCurrency(financial)
}

func (rateApp *RateApp) GetInitialCurrency() ([]entity.Currency, error) {
	return rateApp.dayCurrencyRepository.GetInitialCurrency()
}

func (rateApp *RateApp) SaveRates(financial repository.FinancialRepository) error {
	return rateApp.rateRepository.SaveRates(financial)
}

func (rateApp *RateApp) UpdateRate(id uint64, rate entity.Currency) (entity.Currency, error) {
	return rateApp.rateRepository.UpdateRate(id, rate)
}

func (rateApp *RateApp) DeleteRate(id uint64) error {
	return rateApp.rateRepository.DeleteRate(id)
}

func (rateApp *RateApp) GetRates() ([]entity.Currency, error) {
	return rateApp.rateRepository.GetRates()
}

func (rateApp *RateApp) GetRate(title string) ([]entity.Currency, error) {
	return rateApp.rateRepository.GetRate(title)
}

func (rateApp *RateApp) GetLastRate(title string) (entity.Currency, error) {
	return rateApp.rateRepository.GetLastRate(title)
}
