package repository

import (
	"net/http"
	"server/canal/pkg/domain/entity"
)

type CurrencyLogic interface {
	GetAllLatestCurrencies(r *http.Request) (entity.Description, entity.Currencies, error)
	GetURLCurrencies(r *http.Request) (entity.Description, entity.Currencies, error)
	GetMarkets() (entity.Description, entity.Markets, error)
}
