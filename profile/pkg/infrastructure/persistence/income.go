package persistence

import (
	"context"
	"database/sql"
	"github.com/golang/protobuf/ptypes"
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
		if err = tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{
				"infrastructure": "profile",
				"function":       "GetIncome",
				"action":         "Rollback",
			}).Debug(err)
		}
	}()

	var valueDecimal decimal.Decimal
	switch in.Period {
	case "day":
		err = tx.QueryRow("SELECT value FROM wallet_history WHERE user_id = $1 AND updated_at >= $2::date - interval '1 day' order by updated_at limit 1",
			in.Id,
			time.Now(),
		).Scan(&valueDecimal)
	case "week":
		err = tx.QueryRow("SELECT value FROM wallet_history WHERE user_id = $1 AND updated_at >= $2::date - interval '1 week' order by updated_at limit 1",
			in.Id,
			time.Now(),
		).Scan(&valueDecimal)
	case "month":
		err = tx.QueryRow("SELECT value FROM wallet_history WHERE user_id = $1 AND updated_at >= $2::date - interval '1 month' order by updated_at limit 1",
			in.Id,
			time.Now(),
		).Scan(&valueDecimal)
	case "year":
		err = tx.QueryRow("SELECT value FROM wallet_history WHERE user_id = $1 AND updated_at >= $2::date - interval '1 year' order by updated_at limit 1",
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

func (managerDB *UserDBManager) GetAllIncomePerDay(c context.Context, in *profile.IncomeParameters) (*profile.WalletStates, error) {
	ctx, cancel := context.WithTimeout(c, managerDB.timing)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return &profile.WalletStates{}, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{
				"infrastructure": "profile",
				"function":       "GetAllIncomePerDay",
				"action":         "Rollback",
			}).Debug(err)
		}
	}()

	query := chooseSql(in.Period)
	rows, err := tx.Query(query, in.Id)
	if err != nil {
		return &profile.WalletStates{}, err
	}
	defer rows.Close()

	var out profile.WalletStates
	for rows.Next() {
		var state profile.WalletState
		var updated time.Time

		if err = rows.Scan(&state.Value, &updated); err != nil {
			return &profile.WalletStates{}, err
		}
		if state.UpdatedAt, err = ptypes.TimestampProto(updated); err != nil {
			return &profile.WalletStates{}, err
		}

		out.States = append(out.States, &state)
	}

	if err = rows.Err(); err != nil {
		return &profile.WalletStates{}, err
	}
	if err = tx.Commit(); err != nil {
		return &profile.WalletStates{}, err
	}

	return &out, nil
}

func chooseSql(period string) string {
	switch period {
	case "week":
		return "SELECT value, updated_at FROM wallet_history " +
			"WHERE title = $1 and updated_at between current_date - interval '1 week' and current_date"
	case "month":
		return "SELECT value, updated_at FROM wallet_history " +
			"WHERE title = $1 and updated_at between current_date - interval '1 month' and current_date + interval '1 day'"
	case "year":
		return "SELECT value, updated_at FROM wallet_history " +
			"WHERE title = $1 and updated_at between current_date - interval '1 year' and current_date + interval '1 day'"
	}

	return "SELECT value, updated_at FROM wallet_history " +
		"WHERE title = $1 and updated_at between current_date - interval '1 week' and current_date"
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
			}).Debug(err)
		}
	}()

	err = tx.QueryRow(
		"INSERT INTO wallet_history (user_id, value, updated_at) VALUES ($1, $2, $3)",
		in.Id,
		in.Value,
		time.Now(),
	).Err()
	if err != nil {
		return &profile.Empty{}, err
	}
	if err = tx.Commit(); err != nil {
		return &profile.Empty{}, err
	}

	return &profile.Empty{}, err
}
