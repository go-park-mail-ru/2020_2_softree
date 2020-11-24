package persistence

import (
	"context"
	profile "server/src/profile/pkg/profile/gen"
	"time"
)

func (managerDB *UserDBManager) GetAllPaymentHistory(ctx context.Context, in *profile.UserID) (*profile.AllHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Query(
		"SELECT from_title, to_title, value, amount, updated_at FROM payment_history WHERE user_id = $1",
		in.Id,
	)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var history profile.AllHistory
	for result.Next() {
		var row profile.PaymentHistory
		if err := result.Scan(&row.From, &row.To, &row.Value, &row.Amount, &row.Datetime); err != nil {
			return nil, err
		}

		history.History = append(history.History, &row)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &history, nil
}

func (managerDB *UserDBManager) AddToPaymentHistory(ctx context.Context, in *profile.AddToHistory) (*profile.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	in.Transaction.Datetime.Seconds = time.Now().Unix()
	_, err = tx.Exec(
		"INSERT INTO payment_history (user_id, from_title, to_title, value, amount, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		in.Id,
		in.Transaction.From,
		in.Transaction.To,
		in.Transaction.Value,
		in.Transaction.Amount,
		in.Transaction.Datetime.AsTime(),
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}
