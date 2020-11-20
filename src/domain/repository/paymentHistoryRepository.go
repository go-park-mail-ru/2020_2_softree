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
	GetAllPaymentHistory(uint64) ([]entity.PaymentHistory, error)
}

type paymentHistoryReceiverInterval interface {
	GetIntervalPaymentHistory(uint64, entity.Interval) ([]entity.PaymentHistory, error)
}

type paymentHistoryAddPayment interface {
	AddToPaymentHistory(uint64, entity.PaymentHistory) error
}
