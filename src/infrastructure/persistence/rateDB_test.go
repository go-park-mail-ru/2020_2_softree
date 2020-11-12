package persistence

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"server/src/domain/entity"
	finMock "server/src/infrastructure/mock"
	"testing"
	"time"
)

var testData = map[string]interface{} {
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

//func TestRateDBManager_SaveRatesSuccess(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	require.Equal(t, nil, err)
//	defer db.Close()
//
//	rows := sqlmock.NewRows([]string{"title", "value", "updated_at"})
//	date := time.Now()
//	for name, data := range testData {
//		rows = rows.AddRow(name, data, date)
//	}
//
//	mock.ExpectBegin()
//	for _, name := range ListOfCurrencies {
//		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO history_currency_by_minutes (title, value, updated_at) VALUES ($1, $2, $3)`)).
//			WithArgs(name, testData[name], date).WillReturnResult(sqlmock.NewResult(1, 1))
//	}
//	mock.ExpectCommit()
//
//	ctrl := gomock.NewController(t)
//	mockFinance := finMock.NewFinanceRepositoryForMock(ctrl)
//	mockFinance.EXPECT().GetQuote().Return(testData).Times(len(testData))
//
//	repo := &RateDBManager{DB: db}
//	err = repo.SaveRates(mockFinance)
//	require.NoError(t, err)
//
//	require.Equal(t, nil, mock.ExpectationsWereMet())
//}

func TestRateDBManager_SaveRatesFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO history_currency_by_minutes (title, value, updated_at) VALUES").
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	ctrl := gomock.NewController(t)
	mockFinance := finMock.NewFinanceRepositoryForMock(ctrl)
	mockFinance.EXPECT().GetQuote().Return(testData).Times(len(testData))

	repo := &RateDBManager{DB: db}
	err = repo.SaveRates(mockFinance)

	require.NotEmpty(t, err)
}

func TestRateDBManager_GetRatesSuccess(t *testing.T) {
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

	repo := &RateDBManager{DB: db}
	currencies, err := repo.GetRates()
	require.NoError(t, err)

	for _, curr := range currencies {
		require.EqualValues(t, testData[curr.Title], curr.Value)
		require.EqualValues(t, date, curr.UpdatedAt)
	}
}

func TestRateDBManager_GetRatesFail(t *testing.T) {
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

	repo := &RateDBManager{DB: db}
	_, err = repo.GetRates()
	require.NotEmpty(t, err)
}

func TestRateDBManager_GetRateSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	date := time.Now()
	expected := entity.Currency{Title:"USD", Value: 1.0, UpdatedAt: date}
	rows := sqlmock.NewRows([]string{"value", "updated_at"})
	rows = rows.AddRow(expected.Value, date)

	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT value, updated_at FROM history_currency_by_minutes WHERE").
		WithArgs(expected.Title).
		WillReturnRows(rows)
	mock.ExpectCommit()

	repo := &RateDBManager{DB: db}
	currencies, err := repo.GetRate("USD")
	require.NoError(t, err)

	for _, curr := range currencies {
		require.EqualValues(t, expected.Value, curr.Value)
		require.EqualValues(t, date, curr.UpdatedAt)
	}
}

func TestRateDBManager_GetRateFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.Equal(t, nil, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT value, updated_at FROM history_currency_by_minutes WHERE").
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()

	repo := &RateDBManager{DB: db}
	_, err = repo.GetRate("USD")
	require.NotEmpty(t, err)
}

