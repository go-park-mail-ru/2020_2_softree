package entity

//easyjson:json
type (
	Market struct {
		Base  string `json:"base"`
		Title string `json:"title"`
	}

	Markets []Market
)

func CreateMarkets() Markets {
	return Markets{
		{Base: "USD", Title: "EUR"},
		{Base: "USD", Title: "RUB"},
		{Base: "USD", Title: "JPY"},
		{Base: "USD", Title: "GBP"},
		{Base: "RUB", Title: "GBP"},
		{Base: "RUB", Title: "EUR"},
		{Base: "RUB", Title: "BRL"},
		{Base: "RUB", Title: "ILS"},
		{Base: "RUB", Title: "JPY"},
	}
}
