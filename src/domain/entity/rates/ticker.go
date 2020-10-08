package rates

import (
	"math/rand"
	"time"
)

const (
	randLimit = 150
	interval  = 2
)

func StartTicker() {
	ticker := time.Tick(interval * time.Second)
	rand.Seed(time.Now().UnixNano())
	for range ticker {
		for i := range Currencies {
			newBuy := rand.Float64()*randLimit + 1
			Currencies[i].DoChange(newBuy)
		}
	}
}
