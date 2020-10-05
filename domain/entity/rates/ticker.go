package rates

import (
	"math/rand"
	"time"
)

const randLimit = 150

func StartTicker() {
	ticker := time.Tick(5 * time.Second)
	rand.Seed(time.Now().UnixNano())
	for range ticker {
		for _, obj := range Currencies {
			newBuy := rand.Int() % randLimit + 1
			obj.DoChange(newBuy)
		}
	}
}
