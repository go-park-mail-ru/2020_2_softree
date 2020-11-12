package persistence

import (
	"errors"
	finMock "server/src/infrastructure/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
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

func TestRateDBManager_GetRates(t *testing.T) {

}