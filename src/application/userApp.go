package application

import (
	"server/src/domain/entity"
	"server/src/domain/repository"
)

type UserApp struct {
	userRepository   repository.UserRepository
	paymentHistory   repository.PaymentHistoryRepository
	walletRepository repository.WalletRepository
}

func NewUserApp(
	userRepository repository.UserRepository,
	paymentHistory repository.PaymentHistoryRepository,
	walletRepository repository.WalletRepository) *UserApp {
	return &UserApp{userRepository, paymentHistory, walletRepository}
}

func (userApp *UserApp) SaveUser(u entity.User) (entity.User, error) {
	return userApp.userRepository.SaveUser(u)
}

func (userApp *UserApp) UpdateUserAvatar(id uint64, u entity.User) error {
	return userApp.userRepository.UpdateUserAvatar(id, u)
}

func (userApp *UserApp) UpdateUserPassword(id uint64, u entity.User) error {
	return userApp.userRepository.UpdateUserPassword(id, u)
}

func (userApp *UserApp) DeleteUser(id uint64) error {
	return userApp.userRepository.DeleteUser(id)
}

func (userApp *UserApp) GetUserById(id uint64) (entity.User, error) {
	return userApp.userRepository.GetUserById(id)
}

func (userApp *UserApp) GetUserByLogin(email, password string) (entity.User, error) {
	return userApp.userRepository.GetUserByLogin(email, password)
}

func (userApp *UserApp) GetUserWatchlist(id uint64) ([]entity.Currency, error) {
	return userApp.userRepository.GetUserWatchlist(id)
}

func (userApp *UserApp) CheckExistence(email string) (bool, error) {
	return userApp.userRepository.CheckExistence(email)
}

func (userApp *UserApp) CheckPassword(id uint64, password string) (bool, error) {
	return userApp.userRepository.CheckPassword(id, password)
}

func (userApp *UserApp) GetAllPaymentHistory(id uint64) ([]entity.PaymentHistory, error) {
	return userApp.paymentHistory.GetAllPaymentHistory(id)
}

func (userApp *UserApp) GetIntervalPaymentHistory(id uint64, i entity.Interval) ([]entity.PaymentHistory, error) {
	return userApp.paymentHistory.GetIntervalPaymentHistory(id, i)
}

func (userApp *UserApp) AddToPaymentHistory(id uint64, history entity.PaymentHistory) error {
	return userApp.paymentHistory.AddToPaymentHistory(id, history)
}

func (userApp *UserApp) GetWallets(id uint64) ([]entity.Wallet, error) {
	return userApp.walletRepository.GetWallets(id)
}

func (userApp *UserApp) GetWallet(id uint64, title string) (entity.Wallet, error) {
	return userApp.walletRepository.GetWallet(id, title)
}

func (userApp *UserApp) SetWallet(id uint64, wallet entity.Wallet) error {
	return userApp.walletRepository.SetWallet(id, wallet)
}

func (userApp *UserApp) CheckWallet(id uint64, title string) (bool, error) {
	return userApp.walletRepository.CheckWallet(id, title)
}

func (userApp *UserApp) CreateWallet(id uint64, title string) error {
	return userApp.walletRepository.CreateWallet(id, title)
}
