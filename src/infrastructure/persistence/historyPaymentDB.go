package persistence

import (
	"context"
	"database/sql"
	"server/src/domain/entity"
	"time"

	"github.com/spf13/viper"
)

type PaymentDBManager struct {
	DB *sql.DB
}

func NewPaymentDBManager() (*PaymentDBManager, error) {
	db, err := sql.Open("postgres", viper.GetString("postgres.URL"))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PaymentDBManager{DB: db}, nil
}

func (pm *PaymentDBManager) GetAllPaymentHistory(id uint64) ([]entity.PaymentHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := pm.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Query("SELECT from_title, to_title, value, amount, updated_at FROM payment_history WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	history := make([]entity.PaymentHistory, 0)
	for result.Next() {
		var row entity.PaymentHistory
		if err := result.Scan(&row.From, &row.To, &row.Value, &row.Amount, &row.Datetime); err != nil {
			return nil, err
		}

		history = append(history, row)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return history, nil
}

func (pm *PaymentDBManager) GetIntervalPaymentHistory(id uint64, i entity.Interval) ([]entity.PaymentHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := pm.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Query(
		"SELECT from_title, to_title, value, amount, updated_at FROM payment_history WHERE user_id = $1 AND updated_at BETWEEN $2 AND $3",
		id,
		i.From,
		i.Where,
	)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	history := make([]entity.PaymentHistory, 0)
	for result.Next() {
		var row entity.PaymentHistory
		if err := result.Scan(&row.From, &row.To, &row.Value, &row.Amount, &row.Datetime); err != nil {
			return nil, err
		}

		history = append(history, row)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return history, nil
}

func (pm *PaymentDBManager) AddToPaymentHistory(id uint64, history entity.PaymentHistory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := pm.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	history.Datetime = time.Now()
	_, err = tx.Exec(
		"INSERT INTO payment_history (user_id, from_title, to_title, value, amount, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		id,
		history.From,
		history.To,
		history.Value,
		history.Amount,
		history.Datetime,
	)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
