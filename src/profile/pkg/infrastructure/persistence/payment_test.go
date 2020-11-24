package persistence_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"regexp"
	database "server/src/profile/pkg/infrastructure/persistence"
	profile "server/src/profile/pkg/profile/gen"
	"testing"
)

const (
	userId = 1
	from   = "RUB"
	to     = "USD"
	amount = 200
	value  = 79.9
)

func TestGetAllPaymentHistory_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"from_title", "to_title", "value", "amount", "updated_at"})
	expected := profile.AllHistory{History: []*profile.PaymentHistory{
		{
			From:     from,
			To:       to,
			Amount:   amount,
			Datetime: nil,
			Value:    value,
		},
	}}
	rows = rows.AddRow(
		expected.History[0].From,
		expected.History[0].To,
		expected.History[0].Value,
		expected.History[0].Amount,
		expected.History[0].Datetime,
	)

	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT from_title, to_title, value, amount, updated_at FROM payment_history WHERE").
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
		ExpectQuery("SELECT from_title, to_title, value, amount, updated_at FROM payment_history WHERE").
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
		ExpectExec(regexp.QuoteMeta(`INSERT INTO payment_history (user_id, from_title, to_title, value, amount, updated_at) VALUES`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.AddToPaymentHistory(ctx, &profile.AddToHistory{Id: userId, Transaction: &profile.PaymentHistory{
		From:     from,
		To:       to,
		Amount:   amount,
		Value:    value,
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
		ExpectExec(regexp.QuoteMeta(`INSERT INTO payment_history (user_id, from_title, to_title, value, amount, updated_at) VALUES`)).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.AddToPaymentHistory(ctx, &profile.AddToHistory{Id: userId, Transaction: &profile.PaymentHistory{}})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}
