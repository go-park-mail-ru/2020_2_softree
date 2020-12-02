package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	currency "server/currency/pkg/currency/gen"
	"server/currency/pkg/domain"
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
	db     *sql.DB
	api    domain.FinancialAPI
	timing time.Duration
}

func NewRateDBManager(db *sql.DB, api domain.FinancialAPI) *RateDBManager {
	var timing time.Duration = 5
	if viper.GetDuration("sql.timing") != 0 {
		timing = viper.GetDuration("sql.timing")
	}
	return &RateDBManager{db: db, api: api, timing: timing * time.Second}
}

func (rm *RateDBManager) saveRates(table string, financial domain.FinancialRepository) error {
	currentTime := time.Now()

	query := fmt.Sprintf("INSERT INTO %s (title, value, updated_at) VALUES ($1, $2, $3)", table)

	ctx, cancel := context.WithTimeout(context.Background(), rm.timing)
	defer cancel()

	tx, err := rm.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Println(fmt.Errorf("saveRates: %v", err))
		}
	}()

	for _, name := range ListOfCurrencies {
		quote := financial.GetQuote()[name]
		_, err := tx.Exec(
			query,
			name,
			quote.(float64),
			currentTime,
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

func (rm *RateDBManager) GetAllLatestRates(ctx context.Context, in *currency.Empty) (*currency.Currencies, error) {
	ctx, cancel := context.WithTimeout(ctx, rm.timing)
	defer cancel()

	tx, err := rm.db.BeginTx(ctx, nil)
	if err != nil {
		return &currency.Currencies{}, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{
				"infrastructure": "currency",
				"function":       "GetRates",
				"action":         "Rollback",
			}).Error(err)
		}
	}()

	result, err := tx.Query(
		"SELECT title, value, updated_at FROM history_currency_by_minutes ORDER BY updated_at DESC LIMIT $1",
		LenListOfCurrencies,
	)
	if err != nil {
		return &currency.Currencies{}, err
	}
	defer result.Close()

	var currencies currency.Currencies
	currencies.Rates = make([]*currency.Currency, 0, LenListOfCurrencies)
	for result.Next() {
		var row currency.Currency
		var updatedAt time.Time

		if err := result.Scan(&row.Title, &row.Value, &updatedAt); err != nil {
			return &currency.Currencies{}, err
		}
		if row.UpdatedAt, err = ptypes.TimestampProto(updatedAt); err != nil {
			return &currency.Currencies{}, err
		}

		currencies.Rates = append(currencies.Rates, &row)
	}

	if err = result.Err(); err != nil {
		return &currency.Currencies{}, err
	}
	if err = tx.Commit(); err != nil {
		return &currency.Currencies{}, err
	}

	return &currencies, nil
}

func (rm *RateDBManager) GetAllRatesByTitle(ctx context.Context, in *currency.CurrencyTitle) (*currency.Currencies, error) {
	ctx, cancel := context.WithTimeout(ctx, rm.timing)
	defer cancel()

	tx, err := rm.db.BeginTx(ctx, nil)
	if err != nil {
		return &currency.Currencies{}, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{
				"infrastructure": "currency",
				"function":       "GetRate",
				"action":         "Rollback",
			}).Error(err)
		}
	}()

	result, err := tx.Query("SELECT value, updated_at FROM history_currency_by_minutes WHERE title = $1 ", in.Title)
	if err != nil {
		return &currency.Currencies{}, err
	}
	defer result.Close()

	var currencies currency.Currencies
	currencies.Rates = make([]*currency.Currency, 0, LenListOfCurrencies)
	for result.Next() {
		var row currency.Currency
		var updatedAt time.Time

		row.Title = in.Title
		if err := result.Scan(&row.Value, &updatedAt); err != nil {
			return &currency.Currencies{}, err
		}
		if row.UpdatedAt, err = ptypes.TimestampProto(updatedAt); err != nil {
			return &currency.Currencies{}, err
		}

		currencies.Rates = append(currencies.Rates, &row)
	}
	if err = result.Err(); err != nil {
		return &currency.Currencies{}, err
	}
	if err = tx.Commit(); err != nil {
		return &currency.Currencies{}, err
	}

	return &currencies, nil
}

func (rm *RateDBManager) GetLastRate(ctx context.Context, in *currency.CurrencyTitle) (*currency.Currency, error) {
	ctx, cancel := context.WithTimeout(ctx, rm.timing)
	defer cancel()

	tx, err := rm.db.BeginTx(ctx, nil)
	if err != nil {
		return &currency.Currency{}, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{
				"infrastructure": "currency",
				"function":       "GetLastRate",
				"action":         "Rollback",
			}).Error(err)
		}
	}()
	row := tx.QueryRow(
		"SELECT value, updated_at FROM history_currency_by_minutes WHERE title = $1 ORDER BY updated_at DESC LIMIT 1",
		in.Title,
	)

	result := currency.Currency{Title: in.Title}
	var updatedAt time.Time

	if err = row.Scan(&result.Value, &updatedAt); err != nil {
		return &currency.Currency{}, err
	}
	if result.UpdatedAt, err = ptypes.TimestampProto(updatedAt); err != nil {
		return &currency.Currency{}, err
	}

	if err = tx.Commit(); err != nil {
		return &currency.Currency{}, err
	}

	return &result, nil
}

func (rm *RateDBManager) GetInitialDayCurrency(ctx context.Context, in *currency.Empty) (*currency.InitialDayCurrencies, error) {
	ctx, cancel := context.WithTimeout(ctx, rm.timing)
	defer cancel()

	tx, err := rm.db.BeginTx(ctx, nil)
	if err != nil {
		return &currency.InitialDayCurrencies{}, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{
				"infrastructure": "currency",
				"function":       "GetInitialDayCurrency",
				"action":         "Rollback",
			}).Error(err)
		}
	}()

	result, err := tx.Query(
		"SELECT title, value FROM history_currency_by_minutes ORDER BY updated_at LIMIT $1",
		LenListOfCurrencies,
	)
	if err != nil {
		return &currency.InitialDayCurrencies{}, err
	}
	defer result.Close()

	var currencies currency.InitialDayCurrencies
	currencies.Currencies = make([]*currency.InitialDayCurrency, 0, LenListOfCurrencies)
	for result.Next() {
		var row currency.InitialDayCurrency
		if err := result.Scan(&row.Title, &row.Value); err != nil {
			return &currency.InitialDayCurrencies{}, err
		}

		currencies.Currencies = append(currencies.Currencies, &row)
	}

	if err := result.Err(); err != nil {
		return &currency.InitialDayCurrencies{}, err
	}
	if err = tx.Commit(); err != nil {
		return &currency.InitialDayCurrencies{}, err
	}

	return &currencies, nil
}

func (rm *RateDBManager) truncateTable(table string) error {
	ctx, cancel := context.WithTimeout(context.Background(), rm.timing)
	defer cancel()

	tx, err := rm.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{
				"infrastructure": "currency",
				"function":       "truncateTable",
				"action":         "Rollback",
			}).Error(err)
		}
	}()

	query := fmt.Sprintf("TRUNCATE TABLE %s", table)

	_, err = tx.Exec(query)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
