package persistence

import (
	"context"
	"database/sql"
	currency "server/src/currency/pkg/currency/gen"
	"server/src/currency/pkg/domain"
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

var LenListOfCurrencies = len(ListOfCurrencies)

type RateDBManager struct {
	DB  *sql.DB
	API domain.FinancialAPI
}

func NewRateDBManager(DB *sql.DB, api domain.FinancialAPI) *RateDBManager {
	return &RateDBManager{DB: DB, API: api}
}

func (rm *RateDBManager) saveRates(table string, financial domain.FinancialRepository) error {
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
			"INSERT INTO $4 (title, value, updated_at) VALUES ($1, $2, $3)",
			name,
			quote.(float64),
			currentTime,
			table,
		)

		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (rm *RateDBManager) GetRates(ctx context.Context, in *currency.Empty) (*currency.Currencies, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := rm.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Query(
		"SELECT title, value, updated_at FROM history_currency_by_minutes ORDER BY updated_at DESC LIMIT $1",
		LenListOfCurrencies,
	)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var currencies currency.Currencies
	currencies.Rates = make([]*currency.Currency, 0, LenListOfCurrencies)
	for result.Next() {
		var row currency.Currency
		if err := result.Scan(&row.Title, &row.Value, &row.UpdatedAt); err != nil {
			return nil, err
		}

		currencies.Rates = append(currencies.Rates, &row)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &currencies, nil
}

func (rm *RateDBManager) GetRate(ctx context.Context, in *currency.CurrencyTitle) (*currency.Currencies, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := rm.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Query("SELECT value, updated_at FROM history_currency_by_minutes WHERE title = $1 ", in.Title)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var currencies currency.Currencies
	currencies.Rates = make([]*currency.Currency, 0, LenListOfCurrencies)
	for result.Next() {
		var row currency.Currency
		row.Title = in.Title
		if err := result.Scan(&row.Value, &row.UpdatedAt); err != nil {
			return nil, err
		}

		currencies.Rates = append(currencies.Rates, &row)
	}
	if err = result.Err(); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &currencies, nil
}

func (rm *RateDBManager) GetLastRate(ctx context.Context, in *currency.CurrencyTitle) (*currency.Currency, error) {
	result := currency.Currency{Title: in.Title}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := rm.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	row := tx.QueryRow("SELECT value, updated_at FROM history_currency_by_minutes WHERE title = $1 ORDER BY updated_at DESC LIMIT 1", in.Title)

	if err = row.Scan(&result.Value, &result.UpdatedAt); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &result, nil
}

func (rm *RateDBManager) GetInitialDayCurrency(ctx context.Context, in *currency.Empty) (*currency.InitialDayCurrencies, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := rm.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Query(
		"SELECT title, value FROM history_currency_by_minutes ORDER BY updated_at LIMIT $1",
		LenListOfCurrencies,
	)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var currencies currency.InitialDayCurrencies
	currencies.Currencies = make([]*currency.InitialDayCurrency, 0, LenListOfCurrencies)
	for result.Next() {
		var row currency.InitialDayCurrency
		if err := result.Scan(&row.Title, &row.Value); err != nil {
			return nil, err
		}

		currencies.Currencies = append(currencies.Currencies, &row)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &currencies, nil
}

func (rm *RateDBManager) truncateTable(table string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := rm.DB.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("TRUNCATE TABLE $1", table)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
