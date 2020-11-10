package persistence

import (
	"errors"
	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"server/src/domain/entity"
	"server/src/infrastructure/security"
	"testing"
)

func TestGetUserById_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "email", "password"})
	expected := entity.User{ID: 1, Email: "hound@psina.ru", Password: "long_hashed_string"}
	rows = rows.AddRow(expected.ID, expected.Email, expected.Password)

	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT id, email, password FROM user_trade WHERE").
		WithArgs(uint64(1)).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := &UserDBManager{DB: db}
	row, err := repo.GetUserById(uint64(1))

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, true, reflect.DeepEqual(row, expected))
}

func TestGetUserById_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "email", "password"})
	expected := entity.User{ID: 1, Email: "hound@psina.ru", Password: "long_hashed_string"}
	rows = rows.AddRow(expected.ID, expected.Email, expected.Password)

	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT id, email, password FROM user_trade WHERE").
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := &UserDBManager{DB: db}
	_, err = repo.GetUserById(uint64(1))

	require.NotEmpty(t, err)
}

func TestGetUserByLogin_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "password"})
	expPass, _ := security.MakeShieldedPassword("long_hashed_string")
	expected := entity.User{ID: 1, Email: "hound@psina.ru", Password: expPass}
	rows = rows.AddRow(expected.ID, expected.Password)

	login := "hound@psina.ru"
	password := "long_hashed_string"

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, password FROM user_trade WHERE").WithArgs(login).WillReturnRows(rows)
	mock.ExpectCommit()

	repo := &UserDBManager{DB: db}
	row, err := repo.GetUserByLogin(login, password)

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, true, reflect.DeepEqual(row, expected))
}

func TestGetUserByLogin_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "password"})
	expPass, _ := security.MakeShieldedPassword("long_hashed_string")
	expected := entity.User{ID: 1, Email: "hound@psina.ru", Password: expPass}
	rows = rows.AddRow(expected.ID, expected.Password)

	login := "hound@psina.ru"
	password := "long_hashed_string"

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, password FROM user_trade WHERE").WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := &UserDBManager{DB: db}
	_, err = repo.GetUserByLogin(login, password)

	require.NotEmpty(t, err)
}

func TestGetUserWatchlist_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"base_title", "currency_title"})
	expected := entity.Currency{Base: "USD", Title: "EUR"}
	rows = rows.AddRow(expected.Base, expected.Title)

	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT base_title, currency_title FROM watchlist WHERE").
		WithArgs(uint64(1)).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := &UserDBManager{DB: db}
	row, err := repo.GetUserWatchlist(uint64(1))

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, true, reflect.DeepEqual(row[0], expected))
}

func TestGetUserWatchlist_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"base_title", "currency_title"})
	expected := entity.Currency{Base: "USD", Title: "EUR"}
	rows = rows.AddRow(expected.Base, expected.Title)

	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT base_title, currency_title FROM watchlist WHERE").
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := &UserDBManager{DB: db}
	_, err = repo.GetUserWatchlist(uint64(1))

	require.NotEmpty(t, err)
}

func TestSaveUser_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	login := "hound@psina.ru"
	password := "long_hashed_string"
	expected := entity.User{ID: 1, Email: login, Password: password}

	mock.ExpectBegin()
	mock.
		ExpectQuery("INSERT INTO user_trade (`email`, `password`) VALUES").
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := &UserDBManager{DB: db}
	_, err = repo.SaveUser(expected)

	require.NotEmpty(t, err)
}

func TestDeleteUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM user_trade WHERE").
		WithArgs(uint64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := &UserDBManager{DB: db}
	err = repo.DeleteUser(uint64(1))

	require.Equal(t, nil, err)
}

func TestDeleteUser_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM user_trade WHERE").
		WithArgs(uint64(1)).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := &UserDBManager{DB: db}
	err = repo.DeleteUser(uint64(1))

	require.NotEmpty(t, err)
}
