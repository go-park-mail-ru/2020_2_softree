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
	password, _ := security.MakeShieldedHash(userToSave.Password)
	expectedUser := entity.User{
		ID: 1,
		Email: "yandex@mail.ru",
		Password: password,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := userMock.NewUserRepositoryForMock(ctrl)
	mock.EXPECT().SaveUser(userToSave).Times(1).Return(expectedUser, nil)

	handler := application.UserApp{ur: mock}
	receivedUser, err := handler.ur.SaveUser(userToSave)

	require.NoError(t, err)
	require.Equal(t, expectedUser, receivedUser)
}

func TestUserMemoryRepo_UpdateUserSuccess(t *testing.T) {
	userToSave := entity.User{
		Email: "yandex@mail.ru",
		Password: "password",
	}
	password, _ := security.MakeShieldedHash(userToSave.Password)
	expectedUser := entity.User{
		ID: 1,
		Email: "yandex@mail.ru",
		Password: password,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := userMock.NewUserRepositoryForMock(ctrl)
	mock.EXPECT().SaveUser(userToSave).Times(1).Return(expectedUser, nil)

	handler := application.UserApp{ur: mock}
	receivedUser, err := handler.ur.SaveUser(userToSave)

	require.NoError(t, err)
	require.Equal(t, expectedUser, receivedUser)
}
