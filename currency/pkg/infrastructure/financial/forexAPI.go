package financial

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"server/currency/pkg/domain"
)

type ForexAPI struct {
}

func NewForexAPI() *ForexAPI {
	return &ForexAPI{}
}

type Rate struct {
	Rates map[string]string `json:"rates"`
}

func (f *ForexAPI) GetCurrencies() (domain.FinancialRepository, error) {
	url := "https://api.currencyfreaks.com/latest" +
		"?apikey=" + viper.GetString("api-key.token") +
		"&symbols=USD,RUB,EUR,JPY,GBP,AUD,CAD,CHF,CNY,HKD,NZD,SEK,KRW,SGD,NOK,MXN,INR,ZAR,TRY,BRL,ILS" +
		"&base=USD"

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return &ForexRepo{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return &ForexRepo{}, err
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &ForexRepo{}, err
	}

	var rate Rate
	if err = json.Unmarshal(body, &rate); err != nil {
		return &ForexRepo{}, err
	}

	return convertToForexRepo(rate.Rates), nil
}
