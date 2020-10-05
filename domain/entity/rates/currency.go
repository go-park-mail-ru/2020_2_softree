package rates

type Currency struct {
	title  string
	Buy    int
	Sell   int
	Change float64
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
