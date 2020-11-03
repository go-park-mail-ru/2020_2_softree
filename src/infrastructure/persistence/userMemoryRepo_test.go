package persistence

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"server/src/application"
	"server/src/domain/entity"
	userMock "server/src/infrastructure/mock"
	"server/src/infrastructure/security"
	"testing"
)

func TestUserMemoryRepo_SaveUserSuccess(t *testing.T) {
	userToSave := entity.User{
		Email: "yandex@mail.ru",
		Password: "password",
	}
	password, _ := security.MakeShieldedPassword(userToSave.Password)
	expectedUser := entity.User{
		ID: 1,
		Email: "yandex@mail.ru",
		Password: password,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := userMock.NewUserRepositoryForMock(ctrl)
	mock.EXPECT().SaveUser(userToSave).Times(1).Return(expectedUser, nil)

	handler := application.NewUserApp(mock)
	receivedUser, err := handler.SaveUser(userToSave)

	require.NoError(t, err)
	require.Equal(t, expectedUser, receivedUser)
}
