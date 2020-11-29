package persistence

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"regexp"
	currency "server/currency/pkg/currency/gen"
	mocks "server/currency/pkg/infrastructure/mock"
	"testing"
	"time"
)

var testData = map[string]interface{}{
	"USD": 1.0,
	"RUB": 2.0,
	"EUR": 3.0,
	"JPY": 4.0,
	"GBP": 5.0,
	"AUD": 6.0,
	"CAD": 7.0,
	"CHF": 8.0,
	"CNY": 9.0,
	"HKD": 10.0,
	"NZD": 11.0,
	"SEK": 12.0,
	"KRW": 13.0,
	"SGD": 14.0,
	"NOK": 15.0,
	"MXN": 16.0,
	"INR": 17.0,
	"ZAR": 18.0,
	"TRY": 19.0,
	"BRL": 20.0,
	"ILS": 21.0,
}

func TestRateDBManager_GetRates_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"title", "value", "updated_at"})
	date := time.Now()
	for name, data := range testData {
		rows = rows.AddRow(name, data, date)
	}

	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT title, value, updated_at FROM history_currency_by_minutes ORDER BY updated_at DESC LIMIT").
		WithArgs(uint64(len(testData))).
		WillReturnRows(rows)
	mock.ExpectCommit()

	ctrl := gomock.NewController(t)
	finMock := mocks.NewApiMock(ctrl)

	repo := NewRateDBManager(db, finMock)
	ctx := context.Background()
	currencies, err := repo.GetRates(ctx, nil)
	require.NoError(t, err)

	timestamp, err := ptypes.TimestampProto(date)
	require.NoError(t, err)
	for _, curr := range currencies.Rates {
		require.EqualValues(t, testData[curr.Title], curr.Value)
		require.EqualValues(t, timestamp, curr.UpdatedAt)
	}
}

func TestRateDBManager_GetRates_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"title", "value", "updated_at"})
	date := time.Now()
	for name, data := range testData {
		rows = rows.AddRow(name, data, date)
	}

	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT title, value, updated_at FROM history_currency_by_minutes ORDER BY updated_at DESC LIMIT").
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	ctrl := gomock.NewController(t)
	finMock := mocks.NewApiMock(ctrl)

	repo := NewRateDBManager(db, finMock)
	ctx := context.Background()
	_, err = repo.GetRates(ctx, nil)
	require.NotEmpty(t, err)
}

func TestRateDBManager_GetRate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	date := ptypes.TimestampNow()
	expected := currency.Currency{Title: "USD", Value: 1.0, UpdatedAt: date}
	rows := sqlmock.NewRows([]string{"value", "updated_at"})
	rows = rows.AddRow(expected.Value, date.AsTime())

	mock.ExpectBegin()
	mock.
		ExpectQuery(`SELECT value, updated_at FROM history_currency_by_minutes WHERE`).
		WithArgs(expected.Title).
		WillReturnRows(rows)
	mock.ExpectCommit()

	ctrl := gomock.NewController(t)
	finMock := mocks.NewApiMock(ctrl)

	repo := NewRateDBManager(db, finMock)
	ctx := context.Background()

	title := currency.CurrencyTitle{Title: expected.Title}
	currencies, err := repo.GetRate(ctx, &title)
	require.NoError(t, err)

	for _, curr := range currencies.Rates {
		require.EqualValues(t, expected.Value, curr.Value)
		require.EqualValues(t, date, curr.UpdatedAt)
	}
}

func TestRateDBManager_GetRate_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	date := ptypes.TimestampNow()
	expected := currency.Currency{Title: "USD", Value: 1.0, UpdatedAt: date}
	rows := sqlmock.NewRows([]string{"value", "updated_at"})
	rows.AddRow(expected.Value, date)

	mock.ExpectBegin()
	mock.
		ExpectQuery(`SELECT value, updated_at FROM history_currency_by_minutes WHERE`).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	ctrl := gomock.NewController(t)
	finMock := mocks.NewApiMock(ctrl)

	repo := NewRateDBManager(db, finMock)
	ctx := context.Background()

	title := currency.CurrencyTitle{Title: expected.Title}
	_, err = repo.GetRate(ctx, &title)
	require.NotEmpty(t ,err)
}

func TestRateDBManager_GetLastRate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	date := ptypes.TimestampNow()
	expected := currency.Currency{Title: "USD", Value: 1.0, UpdatedAt: date}
	rows := sqlmock.NewRows([]string{"value", "updated_at"})
	rows = rows.AddRow(expected.Value, date.AsTime())

	mock.ExpectBegin()
	mock.
		ExpectQuery(`SELECT value, updated_at FROM history_currency_by_minutes WHERE title = \$1 ORDER BY updated_at DESC LIMIT 1`).
		WithArgs(expected.Title).
		WillReturnRows(rows)
	mock.ExpectCommit()

	ctrl := gomock.NewController(t)
	finMock := mocks.NewApiMock(ctrl)

	repo := NewRateDBManager(db, finMock)
	ctx := context.Background()

	title := currency.CurrencyTitle{Title: expected.Title}
	res, err := repo.GetLastRate(ctx, &title)
	require.NoError(t, err)

	require.EqualValues(t, res.Title, expected.Title)
	require.EqualValues(t, res.Value, expected.Value)
	require.EqualValues(t, res.UpdatedAt, expected.UpdatedAt)
}

func TestRateDBManager_GetLastRate_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	date := ptypes.TimestampNow()
	expected := currency.Currency{Title: "USD", Value: 1.0, UpdatedAt: date}
	rows := sqlmock.NewRows([]string{"value", "updated_at"})
	rows.AddRow(expected.Value, date)

	mock.ExpectBegin()
	mock.
		ExpectQuery(`SELECT value, updated_at FROM history_currency_by_minutes WHERE title = \$1 ORDER BY updated_at DESC LIMIT 1`).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	ctrl := gomock.NewController(t)
	finMock := mocks.NewApiMock(ctrl)

	repo := NewRateDBManager(db, finMock)
	ctx := context.Background()

	title := currency.CurrencyTitle{Title: expected.Title}
	_, err = repo.GetLastRate(ctx, &title)
	require.NotEmpty(t, err)
}

func TestRateDBManager_GetInitialDayCurrency_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"title", "value"})
	for name, data := range testData {
		rows = rows.AddRow(name, data)
	}

	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT title, value FROM history_currency_by_minutes ORDER BY updated_at LIMIT").
		WithArgs(uint64(len(testData))).
		WillReturnRows(rows)
	mock.ExpectCommit()

	ctrl := gomock.NewController(t)
	finMock := mocks.NewApiMock(ctrl)

	repo := NewRateDBManager(db, finMock)
	ctx := context.Background()
	currencies, err := repo.GetInitialDayCurrency(ctx, nil)
	require.NoError(t, err)

	for _, curr := range currencies.Currencies {
		require.EqualValues(t, testData[curr.Title], curr.Value)
	}
}

func TestRateDBManager_GetInitialDayCurrency_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"title", "value"})
	for name, data := range testData {
		rows = rows.AddRow(name, data)
	}

	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT title, value FROM history_currency_by_minutes ORDER BY updated_at LIMIT").
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	ctrl := gomock.NewController(t)
	finMock := mocks.NewApiMock(ctrl)

	repo := NewRateDBManager(db, finMock)
	ctx := context.Background()
	_, err = repo.GetInitialDayCurrency(ctx, nil)
	require.NotEmpty(t, err)
}

func TestNewRateDBManager_TruncateTable_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	tableName := "history_currency_by_minutes"

	mock.ExpectBegin()
	mock.
		ExpectExec(regexp.QuoteMeta(`TRUNCATE TABLE`)).
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	ctrl := gomock.NewController(t)
	finMock := mocks.NewApiMock(ctrl)

	repo := NewRateDBManager(db, finMock)

	err = repo.truncateTable(tableName)
	require.EqualValues(t, nil, err)
}

func TestNewRateDBManager_TruncateTable_Fail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	tableName := "history_currency_by_minutes"

	mock.ExpectBegin()
	mock.
		ExpectQuery(`TRUNCATE TABLE`).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	ctrl := gomock.NewController(t)
	finMock := mocks.NewApiMock(ctrl)

	repo := NewRateDBManager(db, finMock)

	err = repo.truncateTable(tableName)
	require.NotEmpty(t, err)
}

func TestNewRateDBManager_TruncateTable_FailValidate(t *testing.T) {
	db, _, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	tableName := "table"

	ctrl := gomock.NewController(t)
	finMock := mocks.NewApiMock(ctrl)

	repo := NewRateDBManager(db, finMock)

	err = repo.truncateTable(tableName)
	require.NotEmpty(t, err)
}

func TestRateDBManager_ValidateTable_Success(t *testing.T) {
	for _, name := range tables {
		res := validateTable(name)
		require.EqualValues(t, true, res)
	}
}

func TestRateDBManager_ValidateTable_Fail(t *testing.T) {
	name := "doesNotExist"
	res := validateTable(name)
	require.EqualValues(t, false, res)
}