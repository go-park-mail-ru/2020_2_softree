package repository

import "server/src/domain/entity"

type WalletRepository interface {
	userWalletsReceiver
	userWalletReceiver
	userWalletSet
	userWalletCheck
	userWalletCreateEmpty
}

type userWalletsReceiver interface {
	GetWallets(uint64) ([]entity.Wallet, error)
}

type userWalletReceiver interface {
	GetWallet(uint64, string) (entity.Wallet, error)
}

type userWalletSet interface {
	SetWallet(uint64, entity.Wallet) error
}

type userWalletCheck interface {
	CheckWallet(uint64, string) (bool, error)
}

type userWalletCreateEmpty interface {
	CreateWallet(uint64, string) error
}
