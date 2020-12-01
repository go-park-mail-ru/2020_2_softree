package persistence

import (
	"context"
	"database/sql"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	profile "server/profile/pkg/profile/gen"
)

func (managerDB *UserDBManager) GetIncome(c context.Context, in *profile.IncomeParameters) (*profile.Income, error) {
	ctx, cancel := context.WithTimeout(c, managerDB.timing)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return &profile.Income{}, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{
				"infrastructure": "profile",
				"function":       "CheckExistence",
				"action":         "Rollback",
			}).Error(err)
		}
	}()

	row := tx.QueryRow("SELECT value FROM wallet_history WHERE user_id = $1 AND updated_at = $2", in.Id, in.Period)

	var valueDecimal decimal.Decimal
	if err := row.Scan(&valueDecimal); err != nil {
		if err == sql.ErrNoRows {
			return &profile.Income{Change: 0}, nil
		}
		return &profile.Income{}, err
	}
	if err = tx.Commit(); err != nil {
		return &profile.Income{}, err
	}

	valueFloat, _ := valueDecimal.Float64()
	return &profile.Income{Change: valueFloat}, nil
}
