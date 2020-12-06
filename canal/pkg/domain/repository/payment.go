package repository

import (
	"context"
	"server/canal/pkg/domain/entity"
)

type PaymentLogic interface {
	ReceiveTransactions(ctx context.Context, id int64) (entity.Description, entity.Payments)
	SetPayment(ctx context.Context, payment entity.Payment) entity.Description
}
