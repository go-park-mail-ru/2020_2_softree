package persistence

import (
	"database/sql"
	"server/src/domain/entity"
	"server/src/domain/repository"
	"time"

	_ "github.com/lib/pq"
)

var listOfCurrencies = [...]string{
	"USD",
	"RUB",
	"EUR",
	"JPY",
	"GBP",
	"AUD",
	"CAD",
	"CHF",
	"CNY",
	"HKD",
	"NZD",
	"SEK",
	"KRW",
	"SGD",
	"NOK",
	"MXN",
	"INR",
	"ZAR",
	"TRY",
	"BRL",
	"ILS",
}

type RateDBManager struct {
	DB *sql.DB
}

func (rm *RateDBManager) SaveRates(financial repository.FinancialRepository) error {
	currentTime := time.Now()

	for _, name := range listOfCurrencies {
		quote := financial.GetQuote()[name]
		_, err := rm.DB.Exec(
			"INSERT INTO HistoryCurrencByMinute (`title`, `value`, `updated_at`) VALUES (?, ?, ?)",
			name,
			quote.(float64),
			currentTime,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (rm *RateDBManager) GetRates() ([]entity.Currency, error) {
	result, err := rm.DB.Query(
		"SELECT title, value, updated_at FROM HistoryCurrencByMinute LIMIT ? ORDER BY id DESC",
		len(listOfCurrencies),
	)
	defer result.Close()
	if err != nil {
		return nil, err
	}

	currencies := make([]entity.Currency, len(listOfCurrencies))
	for result.Next() {
		var currency entity.Currency
		if err := result.Scan(currency.Title, currency.Value, currency.UpdatedAt); err != nil {
			return nil, err
		}

		currencies = append(currencies, currency)
	}

	return currencies, nil
}
