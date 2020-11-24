package persistence_test

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"regexp"
	database "server/src/profile/pkg/infrastructure/persistence"
	profile "server/src/profile/pkg/profile/gen"
	"testing"
)

func TestGetWallets_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"title", "value"})
	expected := profile.Wallets{AllWallets: []*profile.Wallet{{Title: from, Value: value}}}
	rows = rows.AddRow(expected.AllWallets[0].Title, expected.AllWallets[0].Value)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT title, value FROM accounts WHERE").WithArgs(userId).WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.GetWallets(ctx, &profile.UserID{Id: userId})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, true, reflect.DeepEqual(row.AllWallets, expected.AllWallets))
}

func TestGetWallets_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT title, value FROM accounts WHERE").
		WithArgs(userId).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.GetWallets(ctx, &profile.UserID{Id: userId})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestGetWallet_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"value"})
	expected := &profile.Wallet{Title: from, Value: value}
	rows = rows.AddRow(expected.Value)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT value FROM accounts WHERE").WithArgs(userId, expected.Title).WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.GetWallet(ctx, &profile.ConcreteWallet{Id: userId, Title: from})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, true, reflect.DeepEqual(row, expected))
}

func TestGetWallet_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT value FROM accounts WHERE").
		WithArgs(userId, from).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.GetWallet(ctx, &profile.ConcreteWallet{Id: userId, Title: from})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestCreateWallet_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO accounts (user_id, title, value) VALUES`)).
		WithArgs(userId, from, decimal.New(0, 0)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.CreateWallet(ctx, &profile.ConcreteWallet{Id: userId, Title: from})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestCreateWallet_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO accounts (user_id, title, value) VALUES`)).
		WithArgs(userId, from, decimal.New(0, 0)).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.CreateWallet(ctx, &profile.ConcreteWallet{Id: userId, Title: from})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestCheckWallet_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"COUNT(user_id)"})
	expected := 0
	rows = rows.AddRow(expected)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(user_id) FROM accounts WHERE`)).
		WithArgs(userId, from).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.CheckWallet(ctx, &profile.ConcreteWallet{Id: userId, Title: from})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, false, row.Existence)
}

func TestCheckWallet_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(user_id) FROM accounts WHERE`)).
		WithArgs(userId, from).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.CheckWallet(ctx, &profile.ConcreteWallet{Id: userId, Title: from})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestSetWallet_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO accounts (user_id, title, value) VALUES`)).
		WithArgs(userId, from, value).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.SetWallet(ctx, &profile.ToSetWallet{Id: userId, NewWallet: &profile.Wallet{Title: from, Value: value}})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestSetWallet_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO accounts (user_id, title, value) VALUES`)).
		WithArgs(userId, from, value).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.SetWallet(ctx, &profile.ToSetWallet{Id: userId, NewWallet: &profile.Wallet{Title: from, Value: value}})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestUpdateWallet_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE accounts SET value = value + $1 WHERE`)).
		WithArgs(decimal.NewFromFloat(value), userId, from).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.UpdateWallet(ctx, &profile.ToSetWallet{Id: userId, NewWallet: &profile.Wallet{Title: from, Value: value}})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestUpdateWallet_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE accounts SET value = value + $1 WHERE`)).
		WithArgs(decimal.NewFromFloat(value), userId, from).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.UpdateWallet(ctx, &profile.ToSetWallet{Id: userId, NewWallet: &profile.Wallet{Title: from, Value: value}})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}
