package persistence

import (
	"database/sql"
	"server/src/domain/entity"
	"server/src/domain/repository"
	"time"

	_ "github.com/lib/pq"
)

type RateDBManager struct {
	DB   *sql.DB
}

func (rm *RateDBManager) SaveRates(financial repository.FinancialRepository) ([]entity.Currency, error) {
	for name, quote := range financial.GetQuote() {
		var rate entity.Currency

		rate.ID = uint64(len(rr.rates) + 1)
		rate.Base = financial.GetBase()
		rate.Title = name
		rate.Value = quote.(float64)

		rates = append(rr.rates, rate)

		result, err := rm.DB.Exec(
			"INSERT INTO HistoryCurrencByMinute (`title`, `value`, `updated_at`) VALUES (?, ?, ?)",
			rate.Title,
			quote.(float64),
			time.Now(),
		)
	}

	return nil
}