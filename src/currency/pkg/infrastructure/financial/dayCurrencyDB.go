package financial

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/shopspring/decimal"
	"server/src/canal/pkg/domain/entity"
	"server/src/currency/pkg/domain"
	persistence2 "server/src/currency/pkg/infrastructure/persistence"
)

type CurrencyManager struct {
	RedisConn redis.Conn
}

func NewCurrencyManager(conn redis.Conn) *CurrencyManager {
	return &CurrencyManager{
		RedisConn: conn,
	}
}

func (sm *CurrencyManager) SaveCurrency(financial domain.FinancialRepository) error {
	for _, name := range persistence2.ListOfCurrencies {
		quote := financial.GetQuote()[name]
		mkey := name
		result, err := redis.String(sm.RedisConn.Do("SET", mkey, quote, "EX", 60*60*24)) // Expires in 24 hours

		if err != nil {
			return err
		}
		if result != "OK" {
			return fmt.Errorf("result not OK")
		}
	}

	return nil
}

func (sm *CurrencyManager) GetInitialCurrency() ([]entity.Currency, error) {
	result := make([]entity.Currency, 0)
	for _, name := range persistence2.ListOfCurrencies {
		var currency entity.Currency
		currency.Title = name
		data, err := redis.Bytes(sm.RedisConn.Do("GET", name))

		if err == redis.ErrNil {
			return []entity.Currency{}, errors.New("no initial value")
		} else if err != nil {
			return []entity.Currency{}, errors.New("redis error during checking initial value")
		}

		strRes := string(data)
		uintRes, parseErr := decimal.NewFromString(strRes)
		if parseErr != nil {
			return []entity.Currency{}, errors.New("internal server error")
		}
		currency.Value = uintRes

		result = append(result, currency)
	}

	return result, nil
}
