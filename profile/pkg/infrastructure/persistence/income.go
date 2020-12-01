package persistence

import (
	"context"
	"database/sql"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	profile "server/profile/pkg/profile/gen"
	"time"
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
				"function":       "GetIncome",
				"action":         "Rollback",
			}).Error(err)
		}
	}()

	var valueDecimal decimal.Decimal
	switch in.Period {
	case "day":
		err = tx.QueryRow("SELECT value FROM wallet_history WHERE user_id = $1 AND updated_at >= DATEADD(day, -1, $2) order by updated_at limit 1",
			in.Id,
			time.Now(),
		).Scan(&valueDecimal)
	case "week":
		err = tx.QueryRow("SELECT value FROM wallet_history WHERE user_id = $1 AND updated_at >= DATEADD(week, -1, $2) order by updated_at limit 1",
			in.Id,
			time.Now(),
		).Scan(&valueDecimal)
	case "month":
		err = tx.QueryRow("SELECT value FROM wallet_history WHERE user_id = $1 AND updated_at >= DATEADD(month, -1, $2) order by updated_at limit 1",
			in.Id,
			time.Now(),
		).Scan(&valueDecimal)
	case "year":
		err = tx.QueryRow("SELECT value FROM wallet_history WHERE user_id = $1 AND updated_at >= DATEADD(year, -1, $2) order by updated_at limit 1",
			in.Id,
			time.Now(),
		).Scan(&valueDecimal)
	case "all_time":
		err = tx.QueryRow("SELECT value FROM wallet_history WHERE user_id = $1 order by updated_at limit 1",
			in.Id,
		).Scan(&valueDecimal)
	}

	if err != nil {
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

func (managerDB *UserDBManager) PutPortfolio(ctx context.Context, in *profile.PortfolioValue) (*profile.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, managerDB.timing)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return &profile.Empty{}, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{
				"infrastructure": "profile",
				"function":       "PutPortfolio",
				"action":         "Rollback",
			}).Error(err)
		}
	}()

	err = tx.
		QueryRow("INSERT INTO wallet_history (user_id, value, updated_at) VALUES ($1, $2, $3)", in.Id, in.Value, time.Now()).Err()
	if err != nil {
		return &profile.Empty{}, err
	}
	if err = tx.Commit(); err != nil {
		return &profile.Empty{}, err
	}

	return &profile.Empty{}, err
}