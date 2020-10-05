package rates

type Currency struct {
	Title  string  `json:"title"`
	Buy    int     `json:"buy"`
	Sell   int     `json:"sell"`
	Change float64 `json:"change"`
}

func (c *Currency) DoChange(newBuy int) {
	if newBuy > c.Buy {
		c.Change = float64(newBuy / c.Buy)
	} else {
		c.Change = -float64(c.Buy / newBuy)
	}

	c.Buy = newBuy
	c.Sell = c.Buy - 1
}

const (
	usdRubPriceBuy  = 78
	usdRubPriceSell = 77
	eurRubPriceBuy  = 88
	eurRubPriceSell = 87
	eurUsdPriceBuy  = 2
	eurUsdPriceSell = 1
	audRubPriceBuy  = 100
	audRubPriceSell = 99
	gelRubPriceBuy  = 20
	gelRubPriceSell = 19
	iskRubPriceBuy  = 120
	iskRubPriceSell = 119
)

var Currencies = []Currency{
	Currency{
		Title:  "USD/RUB",
		Buy:    usdRubPriceBuy,
		Sell:   usdRubPriceSell,
		Change: 0,
	},
	Currency{
		Title:  "EUR/RUB",
		Buy:    eurRubPriceBuy,
		Sell:   eurRubPriceSell,
		Change: 0,
	},
	Currency{
		Title:  "EUR/USD",
		Buy:    eurUsdPriceBuy,
		Sell:   eurUsdPriceSell,
		Change: 0,
	},
	Currency{
		Title:  "AUD/RUB",
		Buy:    audRubPriceBuy,
		Sell:   audRubPriceSell,
		Change: 0,
	},
	Currency{
		Title:  "GEL/RUB",
		Buy:    gelRubPriceBuy,
		Sell:   gelRubPriceSell,
		Change: 0,
	},
	Currency{
		Title:  "ISK/RUB",
		Buy:    iskRubPriceBuy,
		Sell:   iskRubPriceSell,
		Change: 0,
	},
}
