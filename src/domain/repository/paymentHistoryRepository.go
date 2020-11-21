package repository

import (
	"server/src/domain/entity"
)

type PaymentHistoryRepository interface {
	paymentHistoryReceiverAll
	paymentHistoryReceiverInterval
	paymentHistoryAddPayment
}

type paymentHistoryReceiverAll interface {
	GetAllPaymentHistory(int64) ([]entity.PaymentHistory, error)
}

type paymentHistoryReceiverInterval interface {
	GetIntervalPaymentHistory(int64, entity.Interval) ([]entity.PaymentHistory, error)
}

type paymentHistoryAddPayment interface {
	AddToPaymentHistory(int64, entity.PaymentHistory) error
}
