package repository

import (
	"context"
	"net/http"
	"server/canal/pkg/domain/entity"
)

type CurrencyLogic interface {
	GetAllLatestCurrencies(ctx context.Context) (entity.Description, entity.Currencies, error)
	GetURLCurrencies(r *http.Request) (entity.Description, entity.Currencies, error)
	GetMarkets() (entity.Description, entity.Markets, error)
}
