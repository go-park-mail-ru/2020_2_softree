package rates

import (
	"fmt"
	"math/rand"
	"strconv"
)

type Currency struct {
	Title  string  `json:"title"`
	Buy    float64 `json:"buy"`
	Sell   float64 `json:"sell"`
	Change float64 `json:"change"`
}

func (c *Currency) DoChange(newBuy float64) {
	if newBuy > c.Buy {
		c.Change = newBuy / c.Buy
	} else {
		c.Change = -(c.Buy / newBuy)
	}

	c.Buy = newBuy
	c.Sell = c.Buy - rand.Float64()

	formatBuy := fmt.Sprintf("%.2f", c.Buy)
	formatSell := fmt.Sprintf("%.2f", c.Sell)
	formatChange := fmt.Sprintf("%.2f", c.Change)

	c.Buy, _ = strconv.ParseFloat(formatBuy, 64)
	c.Sell, _ = strconv.ParseFloat(formatSell, 64)
	c.Change, _ = strconv.ParseFloat(formatChange, 64)
}

const (
	usdRubPriceBuy  = 78.0
	usdRubPriceSell = 77.0
	eurRubPriceBuy  = 88.0
	eurRubPriceSell = 87.0
	eurUsdPriceBuy  = 2.0
	eurUsdPriceSell = 1.0
	audRubPriceBuy  = 100.0
	audRubPriceSell = 99.0
	gelRubPriceBuy  = 20.0
	gelRubPriceSell = 19.0
	iskRubPriceBuy  = 120.0
	iskRubPriceSell = 119.0
)

var Currencies = []Currency{
	Currency{
		Title:  "USD/RUB",
		Buy:    usdRubPriceBuy,
		Sell:   usdRubPriceSell,
		Change: 0.0,
	},
	Currency{
		Title:  "EUR/RUB",
		Buy:    eurRubPriceBuy,
		Sell:   eurRubPriceSell,
		Change: 0.0,
	},
	Currency{
		Title:  "EUR/USD",
		Buy:    eurUsdPriceBuy,
		Sell:   eurUsdPriceSell,
		Change: 0.0,
	},
	Currency{
		Title:  "AUD/RUB",
		Buy:    audRubPriceBuy,
		Sell:   audRubPriceSell,
		Change: 0.0,
	},
	Currency{
		Title:  "GEL/RUB",
		Buy:    gelRubPriceBuy,
		Sell:   gelRubPriceSell,
		Change: 0.0,
	},
	Currency{
		Title:  "ISK/RUB",
		Buy:    iskRubPriceBuy,
		Sell:   iskRubPriceSell,
		Change: 0.0,
	},
}
