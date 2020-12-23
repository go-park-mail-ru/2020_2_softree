package persistence

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	history_currency_by_minutes = "history_currency_by_minutes"
	history_currency_by_hours   = "history_currency_by_hours"
	history_currency_by_day     = "history_currency_by_day"
)

func (rm *RateDBManager) writeCurrencyDB(table string) {
	finance, err := rm.api.GetCurrencies()
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "writeCurrencyDB", "action": "GetCurrencies"}).Error(err)
		return
	}

	fmt.Println("	GET", "\n---------\n", time.Now(), "\n", finance, "\n---------")

	if err = rm.saveRates(table, finance); err != nil {
		logrus.WithFields(logrus.Fields{"function": "writeCurrencyDB", "action": "saveRates"}).Error(err)
		return
	}

	fmt.Println("	POST", "\n---------\n", time.Now(), "\n", finance, "\n---------")
}

func (rm *RateDBManager) truncate(table string) {
	err := rm.truncateTable(table)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "truncate", "action": "truncateTable"}).Error(err)
		return
	}
}

func (rm *RateDBManager) GetRatesFromApi() {
	task := gocron.NewScheduler(time.UTC)
	defer task.Stop()

	var err error
	if _, err = task.Every(15).
		Minute().StartImmediately().Do(rm.writeCurrencyDB, history_currency_by_minutes); err != nil {
		logrus.WithFields(logrus.Fields{"function": "GetRatesFromApi"}).Error(err)
		return
	}

	if _, err = task.Every(1).
		Hour().StartImmediately().Do(rm.writeCurrencyDB, history_currency_by_hours); err != nil {
		logrus.WithFields(logrus.Fields{"function": "GetRatesFromApi"}).Error(err)
		return
	}

	if _, err = task.Every(1).
		Day().At("00:00").StartImmediately().Do(rm.writeCurrencyDB, history_currency_by_day); err != nil {
		logrus.WithFields(logrus.Fields{"function": "GetRatesFromApi"}).Error(err)
		return
	}

	if _, err = task.Every(1).Day().At("10:00").Do(rm.truncate, history_currency_by_minutes); err != nil {
		logrus.WithFields(logrus.Fields{"function": "GetRatesFromApi"}).Error(err)
		return
	}

	if _, err = task.Every(1).Day().At("00:00").Do(rm.truncate, history_currency_by_minutes); err != nil {
		logrus.WithFields(logrus.Fields{"function": "GetRatesFromApi"}).Error(err)
		return
	}

	if _, err = task.Every(1).Month(1).Do(rm.truncate, history_currency_by_hours); err != nil {
		logrus.WithFields(logrus.Fields{"function": "GetRatesFromApi"}).Error(err)
		return
	}

	<-task.StartAsync()
}
