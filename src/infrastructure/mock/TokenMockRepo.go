package mock

import (
	"github.com/golang/mock/gomock"
	"net/http"
	"reflect"
	"server/src/domain/repository"
)

type TokenRepositoryForMock struct {
	ctrl     *gomock.Controller
	recorder *RecorderTokenMockRepository
}

type RecorderTokenMockRepository struct {
	mock *TokenRepositoryForMock
}

func NewTokenRepositoryForMock(ctrl *gomock.Controller) *TokenRepositoryForMock {
	mock := &TokenRepositoryForMock{ctrl: ctrl}
	mock.recorder = &RecorderTokenMockRepository{mock: mock}
	return mock
}

func (t *TokenRepositoryForMock) EXPECT() *RecorderTokenMockRepository {
	return t.recorder
}

func (t *TokenRepositoryForMock) CreateCookie() (http.Cookie, error) {
	t.ctrl.T.Helper()
	ret := t.ctrl.Call(t, "CreateCookie")
	cookie, _ := ret[0].(http.Cookie)
	err, _ := ret[1].(error)
	return cookie, err
}

func (r *RecorderTokenMockRepository) CreateCookie() *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"CreateCookie",
		reflect.TypeOf((*TokenRepositoryForMock)(nil).CreateCookie),
	)
}

func (t *TokenRepositoryForMock) ExtractData(req *http.Request) (repository.AccessDetails, error) {
	t.ctrl.T.Helper()
	ret := t.ctrl.Call(t, "ExtractData", req)
	details, _ := ret[0].(repository.AccessDetails)
	err, _ := ret[1].(error)
	return details, err
}

func (r *RecorderTokenMockRepository) ExtractData(req interface{}) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(
		r.mock,
		"ExtractData",
		reflect.TypeOf((*TokenRepositoryForMock)(nil).ExtractData),
		req,
	)
}
