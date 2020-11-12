package repository

import "server/src/domain/entity"

type WalletRepository interface {
	userWalletReceiver
	userWalletSet
}

type userWalletReceiver interface {
	GetWallet(uint64) (entity.Wallet, error)
}

type userWalletSet interface {
	SetWallet(uint64, entity.Wallet) error
}
