package repository

import "server/src/domain/entity"

type WalletRepository interface {
	userWalletsReceiver
	userWalletReceiver
	userWalletSet
	userWalletCheck
	userWalletCreateEmpty
	userUpdateWallet
}

type userWalletsReceiver interface {
	GetWallets(int64) ([]entity.Wallet, error)
}

type userWalletReceiver interface {
	GetWallet(int64, string) (entity.Wallet, error)
}

type userWalletSet interface {
	SetWallet(int64, entity.Wallet) error
}

type userWalletCheck interface {
	CheckWallet(int64, string) (bool, error)
}

type userWalletCreateEmpty interface {
	CreateWallet(int64, string) error
}

type userUpdateWallet interface {
	UpdateWallet(int64, entity.Wallet) error
}
