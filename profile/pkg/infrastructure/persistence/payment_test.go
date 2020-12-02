package persistence_test

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"regexp"
	database "server/profile/pkg/infrastructure/persistence"
	profile "server/profile/pkg/profile/gen"
	"testing"
	"time"
)

const (
	userId   = 1
	base     = "RUB"
	currency = "USD"
	amount   = 200
	value    = 79.9
	sell     = "false"
)

func TestGetAllPaymentHistory_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	date := time.Now()
	timestamp, err := ptypes.TimestampProto(date)
	require.NoError(t, err)
	rows := sqlmock.NewRows([]string{"base", "curr", "value", "amount", "sell", "updated_at"})
	expected := profile.AllHistory{History: []*profile.PaymentHistory{
		{
			Base:      base,
			Currency:  currency,
			Amount:    amount,
			UpdatedAt: timestamp,
			Value:     value,
			Sell:      sell,
		},
	}}
	rows = rows.AddRow(
		expected.History[0].Base,
		expected.History[0].Currency,
		expected.History[0].Value,
		expected.History[0].Amount,
		expected.History[0].Sell,
		date,
	)

	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT base, curr, value, amount, sell, updated_at FROM payment_history WHERE").
		WithArgs(userId).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.GetAllPaymentHistory(ctx, &profile.UserID{Id: userId})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, true, reflect.DeepEqual(row.History, expected.History))
}

func TestGetAllPaymentHistory_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT base, curr, value, amount, sell, updated_at FROM payment_history WHERE").
		WithArgs(userId).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.GetAllPaymentHistory(ctx, &profile.UserID{Id: userId})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestAddToPaymentHistory_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.
		ExpectExec(regexp.QuoteMeta(`INSERT INTO payment_history (user_id, base, curr, value, amount, sell, updated_at) VALUES`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.AddToPaymentHistory(ctx, &profile.AddToHistory{Id: userId, Transaction: &profile.PaymentHistory{
		Base:     base,
		Currency: currency,
		Amount:   amount,
		Value:    value,
		Sell:     sell,
	}})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestAddToPaymentHistory_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.
		ExpectExec(regexp.QuoteMeta(`INSERT INTO payment_history (user_id, base, curr, value, amount, sell, updated_at) VALUES`)).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.AddToPaymentHistory(ctx, &profile.AddToHistory{Id: userId, Transaction: &profile.PaymentHistory{}})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}
