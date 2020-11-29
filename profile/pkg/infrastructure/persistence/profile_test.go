package persistence_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"regexp"
	database "server/profile/pkg/infrastructure/persistence"
	profile "server/profile/pkg/profile/gen"
	"testing"
)

const (
	email = "hound@psina.ru"
	password = "password"
	avatar = "base64"
)

func TestGetUserById_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "email", "avatar"})
	expected := profile.PublicUser{Id: userId, Email: email, Avatar: avatar}
	rows = rows.AddRow(expected.Id, expected.Email, expected.Avatar)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, email, avatar FROM user_trade WHERE").
		WithArgs(userId).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.GetUserById(ctx, &profile.UserID{Id: userId})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, true, reflect.DeepEqual(row, &expected))
}

func TestGetUserById_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, email, avatar FROM user_trade WHERE").
		WithArgs(userId).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.GetUserById(ctx, &profile.UserID{Id: userId})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestCheckExistence_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"COUNT(id)"})
	expected := 0
	rows = rows.AddRow(expected)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM user_trade WHERE`)).
		WithArgs(email).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.CheckExistence(ctx, &profile.User{Email: email})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, false, row.Existence)
}

func TestCheckExistence_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM user_trade WHERE`)).
		WithArgs(email).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.CheckExistence(ctx, &profile.User{Email: email})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestGetPassword_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"password"})
	expected := password
	rows = rows.AddRow(expected)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT password FROM user_trade WHERE`)).
		WithArgs(userId).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.GetPassword(ctx, &profile.User{Id: userId, OldPassword: password})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, expected, row.PasswordToCheck)
}

func TestCheckPassword_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT password FROM user_trade WHERE`)).
		WithArgs(userId).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.GetPassword(ctx, &profile.User{Id: userId, OldPassword: password})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestSaveUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id"})
	expected := &profile.PublicUser{Id: userId, Email: email}
	rows = rows.AddRow(expected.Id)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO user_trade (email, password) VALUES`)).
		WithArgs(email, password).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.SaveUser(ctx, &profile.User{Email: email, Password: password})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, false, reflect.DeepEqual(row, &expected))
}

func TestSaveUser_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO user_trade (email, password) VALUES`)).
		WithArgs(email, password).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.SaveUser(ctx, &profile.User{Email: email, Password: password})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestUpdateUserAvatar_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE user_trade SET avatar =`)).
		WithArgs(avatar, userId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.UpdateUserAvatar(ctx, &profile.UpdateFields{Id: userId, User: &profile.User{Avatar: avatar}})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestUpdateUserAvatar_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE user_trade SET avatar =`)).
		WithArgs(avatar, userId).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.UpdateUserAvatar(ctx, &profile.UpdateFields{Id: userId, User: &profile.User{Avatar: avatar}})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestUpdateUserPassword_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE user_trade SET password =`)).
		WithArgs(password, userId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.UpdateUserPassword(ctx, &profile.UpdateFields{Id: userId, User: &profile.User{NewPassword: password}})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestUpdateUserPassword_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE user_trade SET password =`)).
		WithArgs(password, userId).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.UpdateUserPassword(ctx, &profile.UpdateFields{Id: userId, User: &profile.User{NewPassword: password}})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestDeleteUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM user_trade WHERE`)).
		WithArgs(userId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.DeleteUser(ctx, &profile.UserID{Id: userId})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestDeleteUser_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM user_trade WHERE`)).
		WithArgs(userId).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.DeleteUser(ctx, &profile.UserID{Id: userId})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestGetUserByLogin_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	pass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	rows := sqlmock.NewRows([]string{"id", "password", "avatar"})
	expected := profile.PublicUser{Id: userId, Email: email, Avatar: avatar}
	rows = rows.AddRow(expected.Id, pass, expected.Avatar)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, avatar FROM user_trade WHERE`)).
		WithArgs(email).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.GetUserByLogin(ctx, &profile.User{Email: email, Password: password})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, true, reflect.DeepEqual(row, &expected))
}

func TestGetUserByLogin_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, avatar FROM user_trade WHERE`)).
		WithArgs(email).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.GetUserByLogin(ctx, &profile.User{Email: email, Password: password})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestGetUserWatchlist_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"base_title", "currency_title"})
	expected := profile.Currencies{Currencies: []*profile.Currency{
		{
			Base:  from,
			Title: to,
		},
	}}
	rows = rows.AddRow(
		expected.Currencies[0].Base,
		expected.Currencies[0].Title,
	)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT base_title, currency_title FROM watchlist WHERE`)).
		WithArgs(userId).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.GetUserWatchlist(ctx, &profile.UserID{Id: userId})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, true, reflect.DeepEqual(row.Currencies, expected.Currencies))
}

func TestGetUserWatchlistNew_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"base_title", "currency_title"})
	expected := profile.Currencies{Currencies: []*profile.Currency{
		{
			Base:  to,
			Title: from,
		},
	}}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT base_title, currency_title FROM watchlist WHERE`)).
		WithArgs(userId).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.GetUserWatchlist(ctx, &profile.UserID{Id: userId})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, true, reflect.DeepEqual(row.Currencies, expected.Currencies))
}

func TestGetUserWatchlist_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT base_title, currency_title FROM watchlist WHERE`)).
		WithArgs(userId).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.GetUserWatchlist(ctx, &profile.UserID{Id: userId})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}
