package persistence_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"regexp"
	database "server/profile/pkg/infrastructure/persistence"
	profile "server/profile/pkg/profile/gen"
	"testing"
	"time"
)

func TestGetIncome_SuccessAllTime(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"value"})
	rows = rows.AddRow(79.7)

	mock.ExpectBegin()
	mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT value FROM wallet_history WHERE user_id = $1 order by updated_at limit 1`)).
		WithArgs(userId).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.GetIncome(ctx, &profile.IncomeParameters{Id: userId, Period: "all_time"})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, true, reflect.DeepEqual(row.Change, 79.7))
}

func TestGetIncome_SuccessWeek(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"value"})
	rows = rows.AddRow(79.7)

	mock.ExpectBegin()
	mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT value FROM wallet_history WHERE user_id = $1 
			AND updated_at >= current_date - interval '1 week' order by updated_at limit 1`)).
		WithArgs(userId).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.GetIncome(ctx, &profile.IncomeParameters{Id: userId, Period: "week"})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, true, reflect.DeepEqual(row.Change, 79.7))
}

func TestGetIncome_SuccessDay(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"value"})
	rows = rows.AddRow(79.7)

	mock.ExpectBegin()
	mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT value FROM wallet_history WHERE user_id = $1 
			AND updated_at >= current_date - interval '1 day' order by updated_at limit 1`)).
		WithArgs(userId).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.GetIncome(ctx, &profile.IncomeParameters{Id: userId, Period: "day"})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, true, reflect.DeepEqual(row.Change, 79.7))
}

func TestGetIncome_SuccessYear(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"value"})
	rows = rows.AddRow(79.7)

	mock.ExpectBegin()
	mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT value FROM wallet_history WHERE user_id = $1 
			AND updated_at >= current_date - interval '1 year' order by updated_at limit 1`)).
		WithArgs(userId).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.GetIncome(ctx, &profile.IncomeParameters{Id: userId, Period: "year"})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Equal(t, true, reflect.DeepEqual(row.Change, 79.7))
}

func TestGetIncome_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT value FROM wallet_history WHERE user_id = $1 order by updated_at limit 1`)).
		WithArgs(userId).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	_, err = repo.GetIncome(ctx, &profile.IncomeParameters{Id: userId, Period: "all_time"})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
}

func TestGetAllIncomePerDay_SuccessWeek(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"value", "updated_at"})
	rows = rows.AddRow(79.7, time.Now())

	mock.ExpectBegin()
	mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT value, updated_at FROM wallet_history 
			WHERE user_id = $1 and updated_at between current_date - interval '1 week' and current_date`)).
		WithArgs(userId).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.GetAllIncomePerDay(ctx, &profile.IncomeParameters{Id: userId, Period: "week"})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.NotEmpty(t, row.States)
}

func TestGetAllIncomePerDay_SuccessMonth(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"value", "updated_at"})
	rows = rows.AddRow(79.7, time.Now())

	mock.ExpectBegin()
	mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT value, updated_at FROM wallet_history 
			WHERE user_id = $1 and updated_at between current_date - interval '1 month' and current_date`)).
		WithArgs(userId).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.GetAllIncomePerDay(ctx, &profile.IncomeParameters{Id: userId, Period: "month"})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.NotEmpty(t, row.States)
}

func TestGetAllIncomePerDay_SuccessYear(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"value", "updated_at"})
	rows = rows.AddRow(79.7, time.Now())

	mock.ExpectBegin()
	mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT value, updated_at FROM wallet_history 
			WHERE user_id = $1 and updated_at between current_date - interval '1 year' and current_date`)).
		WithArgs(userId).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.GetAllIncomePerDay(ctx, &profile.IncomeParameters{Id: userId, Period: "year"})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.NotEmpty(t, row.States)
}

func TestGetAllIncomePerDay_SuccessEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"value", "updated_at"})
	rows = rows.AddRow(79.7, time.Now())

	mock.ExpectBegin()
	mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT value, updated_at FROM wallet_history 
			WHERE user_id = $1 and updated_at between current_date - interval '1 week' and current_date`)).
		WithArgs(userId).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.GetAllIncomePerDay(ctx, &profile.IncomeParameters{Id: userId})

	require.Equal(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.NotEmpty(t, row.States)
}

func TestGetAllIncomePerDay_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"value", "updated_at"})
	rows = rows.AddRow(79.7, time.Now())

	mock.ExpectBegin()
	mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT value, updated_at FROM wallet_history 
			WHERE user_id = $1 and updated_at between current_date - interval '1 week' and current_date`)).
		WithArgs(userId).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := database.NewUserDBManager(db)
	ctx := context.Background()
	row, err := repo.GetAllIncomePerDay(ctx, &profile.IncomeParameters{Id: userId})

	require.NotEqual(t, nil, err)
	require.Equal(t, nil, mock.ExpectationsWereMet())
	require.Empty(t, row.States)
}
