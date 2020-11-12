package financial

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"server/src/domain/entity"
	"server/src/domain/repository"
	"server/src/infrastructure/persistence"
	"strconv"
)

type CurrencyManager struct {
	RedisConn redis.Conn
}

func NewCurrencyManager(conn redis.Conn) *CurrencyManager {
	return &CurrencyManager{
		RedisConn: conn,
	}
}

func (sm *CurrencyManager) SaveCurrency(financial repository.FinancialRepository) error {
	for _, name := range persistence.ListOfCurrencies {
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
	for _, name := range persistence.ListOfCurrencies {
		var currency entity.Currency
		currency.Title = name
		data, err := redis.Bytes(sm.RedisConn.Do("GET", name))

		if err == redis.ErrNil {
			return []entity.Currency{}, errors.New("no initial value")
		} else if err != nil {
			return []entity.Currency{}, errors.New("redis error during checking initial value")
		}

		strRes := string(data)
		uintRes, parseErr := strconv.ParseFloat(strRes, 64)
		if parseErr != nil {
			return []entity.Currency{}, errors.New("internal server error")
		}
		currency.Value = uintRes

		result = append(result, currency)
	}

	return result, nil
}
