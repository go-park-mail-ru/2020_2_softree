package financial

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/golang/mock/gomock"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/require"
	"server/src/infrastructure/mock"
	"strconv"
	"testing"
)

var testData = map[string]interface{} {
	"USD": 1,
	"RUB": 2,
	"EUR": 3,
	"JPY": 4,
	"GBP": 5,
	"AUD": 6,
	"CAD": 7,
	"CHF": 8,
	"CNY": 9,
	"HKD": 10,
	"NZD": 11,
	"SEK": 12,
	"KRW": 13,
	"SGD": 14,
	"NOK": 15,
	"MXN": 16,
	"INR": 17,
	"ZAR": 18,
	"TRY": 19,
	"BRL": 20,
	"ILS": 21,
}

func TestCurrencyManager_SaveCurrency(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	c, err := redis.Dial("tcp", s.Addr())
	currencyManager := NewCurrencyManager(c)

	ctrl := gomock.NewController(t)
	mockFinance := mock.NewFinanceRepositoryForMock(ctrl)
	mockFinance.EXPECT().GetQuote().Return(testData).Times(len(testData))

	currencyManager.SaveCurrency(mockFinance)

	for name, value := range testData {
		res, err:= s.Get(name)
		require.NoError(t, err)
		strRes, err:= strconv.Atoi(res)
		require.NoError(t, err)
		require.EqualValues(t, value, strRes)
	}
}

func TestCurrencyManager_GetInitialCurrency(t *testing.T) { // Не работает
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	c, err := redis.Dial("tcp", s.Addr())
	currencyManager := NewCurrencyManager(c)

	for name, val := range testData {
		value, _ := val.(int)
		strVal := strconv.Itoa(value)
		require.NoError(t, err)
		s.Set(name, strVal)
	}

	initCurrency, err := currencyManager.GetInitialCurrency()
	require.NoError(t, err)
	for _, currency := range initCurrency {
		require.EqualValues(t, testData[currency.Title], currency.Value)
	}
}
