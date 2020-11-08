package persistence

import (
	"database/sql"
	"fmt"
	"server/src/domain/entity"
	"server/src/domain/repository"
	"time"
	"server/src/infrastructure/config"

	_ "github.com/lib/pq"
)

var ListOfCurrencies = [...]string{
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

func NewRateDBManager() (*RateDBManager, error) {
	/*dsn := config.RateDatabaseConfig.User +
		":" + config.RateDatabaseConfig.Password +
		"@" + config.RateDatabaseConfig.Host +
		"/" + config.RateDatabaseConfig.Schema
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"*/

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.UserDatabaseConfig.Host,
		5432,
		config.UserDatabaseConfig.User,
		config.UserDatabaseConfig.Password,
		config.UserDatabaseConfig.Schema)

	/*dsn := "root:1234@tcp(localhost:3306)/tech?"
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"*/

	db, err := sql.Open("postgres", psqlInfo)

	db.SetMaxOpenConns(10)

	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		return nil, err
	}

	return &RateDBManager{DB: db}, nil
}

func (rm *RateDBManager) SaveRates(financial repository.FinancialRepository) error {
	currentTime := time.Now()

	for _, name := range ListOfCurrencies {
		quote := financial.GetQuote()[name]
		_, err := rm.DB.Exec(
			"INSERT INTO HistoryCurrencByMinute (`title`, `value`, `updated_at`) VALUES ($1, $2, $3)",
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
		"SELECT title, value, updated_at FROM HistoryCurrencByMinute LIMIT $1 ORDER BY id DESC",
		len(ListOfCurrencies),
	)
	defer result.Close()
	if err != nil {
		return nil, err
	}

	currencies := make([]entity.Currency, len(ListOfCurrencies))
	for result.Next() {
		var currency entity.Currency
		if err := result.Scan(&currency.Title, &currency.Value, &currency.UpdatedAt); err != nil {
			return nil, err
		}

		currencies = append(currencies, currency)
	}

	return currencies, nil
}

func (rm *RateDBManager) GetRate(title string) ([]entity.Currency, error) {
	result, err := rm.DB.Query(
		"SELECT value, updated_at FROM HistoryCurrencByMinute WHERE title = $1",
		title,
	)
	defer result.Close()
	if err != nil {
		return nil, err
	}

	currencies := make([]entity.Currency, len(ListOfCurrencies))
	for result.Next() {
		var currency entity.Currency
		currency.Title = title
		if err := result.Scan(&currency.Value, &currency.UpdatedAt); err != nil {
			return nil, err
		}

		currencies = append(currencies, currency)
	}

	return currencies, nil
}

func (rm *RateDBManager) DeleteRate(uint64) error {
	return nil
}

func (rm *RateDBManager) UpdateRate(uint64, entity.Currency) (entity.Currency, error) {
	return entity.Currency{}, nil
}
