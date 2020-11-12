package mock

import (
	"github.com/golang/mock/gomock"
	"reflect"
	"server/src/domain/entity"
)

type UserRepositoryForMock struct {
	ctrl     *gomock.Controller
	recorder *RecorderUserMockRepository
}

type RecorderUserMockRepository struct {
	mock *UserRepositoryForMock
}

func NewUserRepositoryForMock(ctrl *gomock.Controller) *UserRepositoryForMock {
	mock := &UserRepositoryForMock{ctrl: ctrl}
	mock.recorder = &RecorderUserMockRepository{mock: mock}
	return mock
}

func (m *UserRepositoryForMock) EXPECT() *RecorderUserMockRepository {
	return m.recorder
}

func (m *UserRepositoryForMock) SaveUser(u entity.User) (entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUser", u)
	user, _ := ret[0].(entity.User)
	err, _ := ret[1].(error)
	return user, err
}

func (r *RecorderUserMockRepository) SaveUser(u interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"SaveUser",
		reflect.TypeOf((*UserRepositoryForMock)(nil).SaveUser),
		u,
	)
}

func (m *UserRepositoryForMock) UpdateUserAvatar(id uint64, u entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserAvatar", id, u)
	err, _ := ret[0].(error)
	return err
}

func (r *RecorderUserMockRepository) UpdateUserAvatar(id, u interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"UpdateUserAvatar",
		reflect.TypeOf((*UserRepositoryForMock)(nil).UpdateUserAvatar),
		id,
		u,
	)
}

func (m *UserRepositoryForMock) UpdateUserPassword(id uint64, u entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPassword", id, u)
	err, _ := ret[0].(error)
	return err
}

func (r *RecorderUserMockRepository) UpdateUserPassword(id, u interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"UpdateUserPassword",
		reflect.TypeOf((*UserRepositoryForMock)(nil).UpdateUserPassword),
		id,
		u,
	)
}

func (m *UserRepositoryForMock) DeleteUser(id uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", id)
	err, _ := ret[0].(error)
	return err
}

func (r *RecorderUserMockRepository) DeleteUser(id interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"DeleteUser",
		reflect.TypeOf((*UserRepositoryForMock)(nil).DeleteUser),
		id,
	)
}

func (m *UserRepositoryForMock) GetUserById(id uint64) (entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", id)
	user, _ := ret[0].(entity.User)
	err, _ := ret[1].(error)
	return user, err
}

func (r *RecorderUserMockRepository) GetUserById(id interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"GetUserById",
		reflect.TypeOf((*UserRepositoryForMock)(nil).GetUserById),
		id,
	)
}

func (m *UserRepositoryForMock) GetUserByLogin(email, password string) (entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", email, password)
	user, _ := ret[0].(entity.User)
	err, _ := ret[1].(error)
	return user, err
}

func (r *RecorderUserMockRepository) GetUserByLogin(email, password interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"GetUserByLogin",
		reflect.TypeOf((*UserRepositoryForMock)(nil).GetUserByLogin),
		email,
		password,
	)
}

func (m *UserRepositoryForMock) CheckExistence(email string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckExistence", email)
	user, _ := ret[0].(bool)
	err, _ := ret[1].(error)
	return user, err
}

func (r *RecorderUserMockRepository) CheckExistence(email interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"CheckExistence",
		reflect.TypeOf((*UserRepositoryForMock)(nil).CheckExistence),
		email,
	)
}

func (m *UserRepositoryForMock) GetUserWatchlist(id uint64) ([]entity.Currency, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUser", id)
	currencies, _ := ret[0].([]entity.Currency)
	err, _ := ret[1].(error)
	return currencies, err
}

func (r *RecorderUserMockRepository) GetUserWatchlist(id interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"SaveUser",
		reflect.TypeOf((*UserRepositoryForMock)(nil).GetUserWatchlist),
		id,
	)
}
