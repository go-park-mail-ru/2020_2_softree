package userInteraction

import (
	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"server/src/domain/entity"
	"server/src/infrastructure/security"
	"testing"
)

func TestUserDBManager_GetUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "email", "password"})
	expected := entity.User{ID: 1, Email: "hound@psina.ru", Password: "long_hashed_string"}
	rows = rows.AddRow(expected.ID, expected.Email, expected.Password)

	mock.ExpectQuery("SELECT id, email, password FROM user_trade WHERE").WithArgs(uint64(1)).WillReturnRows(rows)

	repo := &UserDBManager{DB: db}
	row, err := repo.GetUserById(uint64(1))

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, true, reflect.DeepEqual(row, expected))
}

func TestUserDBManager_GetUserByLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "password"})
	expPass, _ := security.MakeShieldedPassword("long_hashed_string")
	expected := entity.User{ID: 1, Email: "hound@psina.ru", Password: expPass}
	rows = rows.AddRow(expected.ID, expected.Password)

	login := "hound@psina.ru"
	password := "long_hashed_string"
	mock.ExpectQuery("SELECT id, password FROM user_trade WHERE").WithArgs(login).WillReturnRows(rows)

	repo := &UserDBManager{DB: db}
	row, err := repo.GetUserByLogin(login, password)

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, true, reflect.DeepEqual(row, expected))
}
