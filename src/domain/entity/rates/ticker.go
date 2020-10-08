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
	ticker := time.NewTicker(interval * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		for i := range Currencies {
			newBuy := rand.Float64()*randLimit + 1
			Currencies[i].DoChange(newBuy)
		}
	}
}
