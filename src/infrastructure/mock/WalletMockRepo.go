package mock

import (
	"github.com/golang/mock/gomock"
	"reflect"
	"server/src/domain/entity"
)

type WalletRepositoryForMock struct {
	ctrl     *gomock.Controller
	recorder *RecorderWalletMockRepository
}

type RecorderWalletMockRepository struct {
	mock *WalletRepositoryForMock
}

func NewWalletRepositoryForMock(ctrl *gomock.Controller) *WalletRepositoryForMock {
	mock := &WalletRepositoryForMock{ctrl: ctrl}
	mock.recorder = &RecorderWalletMockRepository{mock: mock}
	return mock
}

func (w *WalletRepositoryForMock) EXPECT() *RecorderWalletMockRepository {
	return w.recorder
}

func (w *WalletRepositoryForMock) GetWallet(id uint64) ([]entity.Wallet, error) {
	w.ctrl.T.Helper()
	ret := w.ctrl.Call(w, "GetWallet", id)
	wallet, _ := ret[0].([]entity.Wallet)
	err, _ := ret[1].(error)
	return wallet, err
}

func (r *RecorderWalletMockRepository) GetWallet(id interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"GetWallet",
		reflect.TypeOf((*WalletRepositoryForMock)(nil).GetWallet),
		id,
	)
}

func (w *WalletRepositoryForMock) SetWallet(id uint64, wallet entity.Wallet) error {
	w.ctrl.T.Helper()
	ret := w.ctrl.Call(w, "GetWallet", id, wallet)
	err, _ := ret[0].(error)
	return err
}

func (r *RecorderWalletMockRepository) SetWallet(id, wallet interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"SetWallet",
		reflect.TypeOf((*WalletRepositoryForMock)(nil).SetWallet),
		id,
		wallet,
	)
}
