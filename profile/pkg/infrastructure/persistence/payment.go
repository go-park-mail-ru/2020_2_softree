package persistence

import (
	"context"
	"github.com/sirupsen/logrus"
	profile "server/profile/pkg/profile/gen"
	"time"

	"github.com/golang/protobuf/ptypes"
)

func (managerDB *UserDBManager) GetAllPaymentHistory(c context.Context, in *profile.UserID) (*profile.AllHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
			}).Error(err)
		}
	}()

	result, err := tx.Query(
		"SELECT from_title, to_title, value, amount, updated_at FROM payment_history WHERE user_id = $1",
		in.Id,
	)
	if err != nil {
		return &profile.AllHistory{}, err
	}
	defer result.Close()

	var history profile.AllHistory
	for result.Next() {
		var row profile.PaymentHistory
		var updatedAt time.Time

		if err := result.Scan(&row.From, &row.To, &row.Value, &row.Amount, &updatedAt); err != nil {
			return &profile.AllHistory{}, err
		}
		if row.UpdatedAt, err = ptypes.TimestampProto(updatedAt); err != nil {
			return &profile.AllHistory{}, err
		}

		history.History = append(history.History, &row)
	}

	if err := result.Err(); err != nil {
		return &profile.AllHistory{}, err
	}
	if err = tx.Commit(); err != nil {
		return &profile.AllHistory{}, err
	}

	return &history, nil
}

func (managerDB *UserDBManager) AddToPaymentHistory(c context.Context, in *profile.AddToHistory) (*profile.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
			}).Error(err)
		}
	}()

	in.Transaction.UpdatedAt = ptypes.TimestampNow()
	_, err = tx.Exec(
		"INSERT INTO payment_history (user_id, from_title, to_title, value, amount, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		in.Id,
		in.Transaction.From,
		in.Transaction.To,
		in.Transaction.Value,
		in.Transaction.Amount,
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
