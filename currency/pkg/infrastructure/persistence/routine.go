package persistence

import (
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"server/currency/pkg/domain"
	"time"
)

const (
	history_currency_by_minutes = "history_currency_by_minutes"
	history_currency_by_hours   = "history_currency_by_hours"
	history_currency_by_day     = "history_currency_by_day"
)

func (rm *RateDBManager) writeCurrencyDB(table string, finance domain.FinancialRepository) {
	err := rm.saveRates(table, finance)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "writeCurrencyDB",
			"action":   "saveRates",
		}).Error(err)
		return
	}
}

func (rm *RateDBManager) truncate(table string) {
	err := rm.truncateTable(table)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "truncate",
			"action":   "truncateTable",
		}).Error(err)
		return
	}
}

func (rm *RateDBManager) GetRatesFromApi() {
	finance := rm.api.GetCurrencies()
	task := gocron.NewScheduler(time.UTC)
	defer task.Stop()

	if _, err := task.Every(1).Minute().Do(rm.writeCurrencyDB, history_currency_by_minutes, finance); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetRatesFromApi",
		}).Error(err)
		return
	}
	if _, err := task.Every(1).Hour().Do(rm.writeCurrencyDB, history_currency_by_hours, finance); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetRatesFromApi",
		}).Error(err)
		return
	}
	if _, err := task.Every(1).Day().At("00:00").Do(rm.writeCurrencyDB, history_currency_by_day, finance); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetRatesFromApi",
		}).Error(err)
		return
	}
	if _, err := task.Every(1).Day().At("00:00").Do(rm.truncate, history_currency_by_minutes); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetRatesFromApi",
		}).Error(err)
		return
	}
	if _, err := task.Every(1).Month(1).Do(rm.truncate, history_currency_by_hours); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetRatesFromApi",
		}).Error(err)
		return
	}

	<-task.StartAsync()
}
