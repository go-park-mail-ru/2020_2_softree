package financial

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"server/src/domain/repository"
	"server/src/infrastructure/auth"
	"server/src/infrastructure/persistence"
)

func (sm *auth.SessionManager) SaveCurrency(financial repository.FinancialRepository) error {
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
