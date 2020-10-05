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
		for i, _ := range Currencies {
			newBuy := rand.Int() % randLimit + 1
			Currencies[i].DoChange(newBuy)
		}
	}
}
