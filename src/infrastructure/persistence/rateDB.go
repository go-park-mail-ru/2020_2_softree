package persistence

import (
	"database/sql"
	"server/src/domain/entity"
	"server/src/domain/repository"

	_ "github.com/lib/pq"
)

type RateDBManager struct {
	DB   *sql.DB
}

func (rm *RateDBManager) SaveRates(financial repository.FinancialRepository) error {
	for name, quote := range financial.GetQuote() {
		var rate entity.Currency

		rate.ID = uint64(len(rr.rates) + 1)
		rate.Base = financial.GetBase()
		rate.Title = name
		rate.Value = quote.(float64)

		rr.rates = append(rr.rates, rate)

		result, err := rm.DB.Exec(
			"INSERT INTO user (`email`, `password`) VALUES (?, ?)",
			user.Email,
			password,
		)
	}

	return nil
}