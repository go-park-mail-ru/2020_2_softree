package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"server/src/domain/entity"
	"server/src/domain/repository"
	"server/src/infrastructure/config"
	"time"
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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.UserDatabaseConfig.Host,
		5432,
		config.UserDatabaseConfig.User,
		config.UserDatabaseConfig.Password,
		config.UserDatabaseConfig.Schema)

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := rm.DB.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, name := range ListOfCurrencies {
		quote := financial.GetQuote()[name]
		_, err := tx.Exec(
			"INSERT INTO HistoryCurrencByMinute (`title`, `value`, `updated_at`) VALUES ($1, $2, $3)",
			name,
			quote.(float64),
			currentTime,
		)

		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (rm *RateDBManager) GetRates() ([]entity.Currency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := rm.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Query(
		"SELECT title, value, updated_at FROM HistoryCurrencByMinute LIMIT $1 ORDER BY id DESC",
		len(ListOfCurrencies),
	)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	currencies := make([]entity.Currency, len(ListOfCurrencies))
	for result.Next() {
		var currency entity.Currency
		if err := result.Scan(&currency.Title, &currency.Value, &currency.UpdatedAt); err != nil {
			return nil, err
		}

		currencies = append(currencies, currency)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return currencies, nil
}

func (rm *RateDBManager) GetRate(title string) ([]entity.Currency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := rm.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Query("SELECT value, updated_at FROM HistoryCurrencByMinute WHERE title = $1", title)
	defer result.Close()
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
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
