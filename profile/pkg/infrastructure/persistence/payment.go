package persistence

import (
	"context"
	"github.com/sirupsen/logrus"
	profile "server/profile/pkg/profile/gen"
	"time"

	"github.com/golang/protobuf/ptypes"
)

func (managerDB *UserDBManager) GetAllPaymentHistory(c context.Context, in *profile.UserID) (*profile.AllHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), managerDB.timing)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return &profile.AllHistory{}, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{
				"infrastructure": "profile",
				"function":       "GetAllPaymentHistory",
				"action":         "Rollback",
			}).Debug(err)
		}
	}()

	result, err := tx.Query(
		"SELECT base, curr, value, amount, sell, updated_at FROM payment_history WHERE user_id = $1 order by updated_at desc",
		in.Id,
	)
	if err != nil {
		return &profile.AllHistory{}, err
	}
	defer result.Close()

	var history profile.AllHistory
	for result.Next() {
		var pay profile.PaymentHistory
		var updatedAt time.Time

		if err := result.Scan(&pay.Base, &pay.Currency, &pay.Value, &pay.Amount, &pay.Sell, &updatedAt); err != nil {
			return &profile.AllHistory{}, err
		}
		if pay.UpdatedAt, err = ptypes.TimestampProto(updatedAt); err != nil {
			return &profile.AllHistory{}, err
		}

		history.History = append(history.History, &pay)
	}

	if err = result.Err(); err != nil {
		return &profile.AllHistory{}, err
	}
	if err = tx.Commit(); err != nil {
		return &profile.AllHistory{}, err
	}

	return &history, nil
}

func (managerDB *UserDBManager) AddToPaymentHistory(c context.Context, in *profile.AddToHistory) (*profile.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), managerDB.timing)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return &profile.Empty{}, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{
				"infrastructure": "profile",
				"function":       "AddToPaymentHistory",
				"action":         "Rollback",
			}).Debug(err)
		}
	}()

	in.Transaction.UpdatedAt = ptypes.TimestampNow()
	query := "INSERT INTO payment_history (user_id, base, curr, value, amount, sell, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err = tx.Exec(
		query,
		in.Id,
		in.Transaction.Base,
		in.Transaction.Currency,
		in.Transaction.Value,
		in.Transaction.Amount,
		in.Transaction.Sell,
		in.Transaction.UpdatedAt.AsTime(),
	)
	if err != nil {
		return &profile.Empty{}, err
	}
	if err = tx.Commit(); err != nil {
		return &profile.Empty{}, err
	}

	return &profile.Empty{}, nil
}
