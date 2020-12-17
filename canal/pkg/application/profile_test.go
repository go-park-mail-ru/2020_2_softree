package application_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"reflect"
	"server/canal/pkg/application"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/infrastructure/mock"
	profile "server/profile/pkg/infrastructure/mock"
	"server/profile/pkg/profile/gen"
	"testing"
)

const (
	id          = int64(1)
	email       = "hound@psina.ru"
	oldPassword = "old"
	newPassword = "new"
	avatar      = "base64"
	currValue   = 79.7
	baseValue   = 1.0
	amount      = 1000.0
	base        = "RUB"
	curr        = "USD"
	sell        = false
)

func TestReceiveUser_Success(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createReceiveUserSuccess(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.ReceiveUser(ctx, id)

	require.NoError(t, err)
	require.Empty(t, desc)
	require.NotEmpty(t, out)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(out))
}

func createReceiveUserSuccess(t *testing.T, ctx context.Context) (*application.ProfileApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		GetUserById(ctx, &gen.UserID{Id: id}).
		Return(&gen.PublicUser{Id: id, Email: email, Avatar: avatar}, nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewProfileApp(profileService, securityService), ctrl
}

func TestReceiveUser_Fail(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createReceiveUserFail(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.ReceiveUser(ctx, id)

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Empty(t, out)
	require.Equal(t, http.StatusInternalServerError, desc.Status)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(out))
}

func createReceiveUserFail(t *testing.T, ctx context.Context) (*application.ProfileApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		GetUserById(ctx, &gen.UserID{Id: id}).
		Return(&gen.PublicUser{}, errors.New("error"))

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewProfileApp(profileService, securityService), ctrl
}

func TestReceiveWatchlist_Success(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createReceiveWatchlistSuccess(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.ReceiveWatchlist(ctx, id)

	require.NoError(t, err)
	require.Empty(t, desc)
	require.NotEmpty(t, out)
	require.Equal(t, reflect.TypeOf(entity.Currencies{}), reflect.TypeOf(out))
}

func createReceiveWatchlistSuccess(t *testing.T, ctx context.Context) (*application.ProfileApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		GetUserWatchlist(ctx, &gen.UserID{Id: id}).
		Return(createWatchlist(), nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewProfileApp(profileService, securityService), ctrl
}

func TestReceiveWatchlist_Fail(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createReceiveWatchlistFail(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.ReceiveWatchlist(ctx, id)

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Empty(t, out)
	require.Equal(t, http.StatusInternalServerError, desc.Status)
	require.Equal(t, reflect.TypeOf(entity.Currencies{}), reflect.TypeOf(out))
}

func createReceiveWatchlistFail(t *testing.T, ctx context.Context) (*application.ProfileApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		GetUserWatchlist(ctx, &gen.UserID{Id: id}).
		Return(&gen.Currencies{}, errors.New("error"))

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewProfileApp(profileService, securityService), ctrl
}

func TestUpdateAvatar_Success(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createUpdateAvatarSuccess(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.UpdateAvatar(ctx, entity.User{Id: id, Avatar: avatar})

	require.NoError(t, err)
	require.Empty(t, desc)
	require.NotEmpty(t, out)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(out))
}

func createUpdateAvatarSuccess(t *testing.T, ctx context.Context) (*application.ProfileApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		UpdateUserAvatar(ctx, &gen.UpdateFields{Id: id, User: &gen.User{Id: id, Avatar: avatar}}).
		Return(&gen.Empty{}, nil)
	profileService.EXPECT().
		GetUserById(ctx, &gen.UserID{Id: id}).
		Return(&gen.PublicUser{Id: id, Email: email, Avatar: avatar}, nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewProfileApp(profileService, securityService), ctrl
}

func TestUpdateAvatar_FailUpdateAvatar(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createUpdateAvatarFailUpdateAvatar(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.UpdateAvatar(ctx, entity.User{Id: id, Avatar: avatar})

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Empty(t, out)
	require.Equal(t, http.StatusInternalServerError, desc.Status)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(out))
}

func createUpdateAvatarFailUpdateAvatar(t *testing.T, ctx context.Context) (*application.ProfileApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		UpdateUserAvatar(ctx, &gen.UpdateFields{Id: id, User: &gen.User{Id: id, Avatar: avatar}}).
		Return(&gen.Empty{}, errors.New("error"))

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewProfileApp(profileService, securityService), ctrl
}

func TestUpdateAvatar_FailGetById(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createUpdateAvatarFailGetById(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.UpdateAvatar(ctx, entity.User{Id: id, Avatar: avatar})

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Empty(t, out)
	require.Equal(t, http.StatusInternalServerError, desc.Status)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(out))
}

func createUpdateAvatarFailGetById(t *testing.T, ctx context.Context) (*application.ProfileApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		UpdateUserAvatar(ctx, &gen.UpdateFields{Id: id, User: &gen.User{Id: id, Avatar: avatar}}).
		Return(&gen.Empty{}, nil)
	profileService.EXPECT().
		GetUserById(ctx, &gen.UserID{Id: id}).
		Return(&gen.PublicUser{}, errors.New("error"))

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewProfileApp(profileService, securityService), ctrl
}

func TestUpdateAvatar_FailValidation(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createUpdateAvatarFailValidation(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.UpdateAvatar(ctx, entity.User{Id: id})

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Empty(t, out)
	require.Equal(t, http.StatusBadRequest, desc.Status)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(out))
}

func createUpdateAvatarFailValidation(t *testing.T, ctx context.Context) (*application.ProfileApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewProfileApp(profileService, securityService), ctrl
}

func TestUpdatePassword_Success(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createUpdatePasswordSuccess(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.UpdatePassword(ctx, entity.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword})

	require.NoError(t, err)
	require.Empty(t, desc)
	require.NotEmpty(t, out)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(out))
}

func createUpdatePasswordSuccess(t *testing.T, ctx context.Context) (*application.ProfileApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		GetPassword(ctx, &gen.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword}).
		Return(&gen.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword, PasswordToCheck: oldPassword}, nil)
	profileService.EXPECT().
		UpdateUserPassword(ctx, &gen.UpdateFields{Id: id, User: &gen.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword, PasswordToCheck: oldPassword}}).
		Return(&gen.Empty{}, nil)
	profileService.EXPECT().
		GetUserById(ctx, &gen.UserID{Id: id}).
		Return(&gen.PublicUser{Id: id, Email: email, Avatar: avatar}, nil)

	securityService := mock.NewSecurityMock(ctrl)
	securityService.EXPECT().
		CheckPassword(oldPassword, oldPassword).
		Return(true)
	securityService.EXPECT().
		MakeShieldedPassword(newPassword).
		Return(newPassword, nil)

	return application.NewProfileApp(profileService, securityService), ctrl
}

func TestUpdatePassword_FailValidationV1(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createUpdatePasswordFailValidationV1(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.UpdatePassword(ctx, entity.User{Id: id})

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Empty(t, out)
	require.Equal(t, http.StatusBadRequest, desc.Status)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(out))
}

func createUpdatePasswordFailValidationV1(t *testing.T, ctx context.Context) (*application.ProfileApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewProfileApp(profileService, securityService), ctrl
}

func TestUpdatePassword_FailValidationV2(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createUpdatePasswordFailValidationV2(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.UpdatePassword(ctx, entity.User{Id: id, OldPassword: "cd cd", NewPassword: "cd fv"})

	require.NoError(t, err)
	require.NotEmpty(t, desc.ErrorJSON)
	require.Empty(t, out)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(out))
}

func createUpdatePasswordFailValidationV2(t *testing.T, ctx context.Context) (*application.ProfileApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewProfileApp(profileService, securityService), ctrl
}

func TestUpdatePassword_FailGetPassword(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createUpdatePasswordFailGetPassword(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.UpdatePassword(ctx, entity.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword})

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Empty(t, out)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(out))
}

func createUpdatePasswordFailGetPassword(t *testing.T, ctx context.Context) (*application.ProfileApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		GetPassword(ctx, &gen.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword}).
		Return(&gen.User{}, errors.New("error"))

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewProfileApp(profileService, securityService), ctrl
}

func TestUpdatePassword_FailCheckPassword(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createUpdatePasswordFailCheckPassword(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.UpdatePassword(ctx, entity.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword})

	require.NoError(t, err)
	require.NotEmpty(t, desc)
	require.Empty(t, out)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(out))
}

func createUpdatePasswordFailCheckPassword(t *testing.T, ctx context.Context) (*application.ProfileApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		GetPassword(ctx, &gen.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword}).
		Return(&gen.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword, PasswordToCheck: oldPassword}, nil)

	securityService := mock.NewSecurityMock(ctrl)
	securityService.EXPECT().
		CheckPassword(oldPassword, oldPassword).
		Return(false)

	return application.NewProfileApp(profileService, securityService), ctrl
}

func TestUpdatePassword_FailMakeHash(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createUpdatePasswordFailMakeHash(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.UpdatePassword(ctx, entity.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword})

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Empty(t, out)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(out))
}

func createUpdatePasswordFailMakeHash(t *testing.T, ctx context.Context) (*application.ProfileApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		GetPassword(ctx, &gen.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword}).
		Return(&gen.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword, PasswordToCheck: oldPassword}, nil)

	securityService := mock.NewSecurityMock(ctrl)
	securityService.EXPECT().
		CheckPassword(oldPassword, oldPassword).
		Return(true)
	securityService.EXPECT().
		MakeShieldedPassword(newPassword).
		Return("", errors.New("error"))

	return application.NewProfileApp(profileService, securityService), ctrl
}

func TestUpdatePassword_FailUpdatePassword(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createUpdatePasswordFailUpdatePassword(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.UpdatePassword(ctx, entity.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword})

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Empty(t, out)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(out))
}

func createUpdatePasswordFailUpdatePassword(t *testing.T, ctx context.Context) (*application.ProfileApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		GetPassword(ctx, &gen.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword}).
		Return(&gen.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword, PasswordToCheck: oldPassword}, nil)
	profileService.EXPECT().
		UpdateUserPassword(ctx, &gen.UpdateFields{Id: id, User: &gen.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword, PasswordToCheck: oldPassword}}).
		Return(&gen.Empty{}, errors.New("error"))

	securityService := mock.NewSecurityMock(ctrl)
	securityService.EXPECT().
		CheckPassword(oldPassword, oldPassword).
		Return(true)
	securityService.EXPECT().
		MakeShieldedPassword(newPassword).
		Return(newPassword, nil)

	return application.NewProfileApp(profileService, securityService), ctrl
}

func TestUpdatePassword_FailGetById(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createUpdatePasswordFailGetById(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.UpdatePassword(ctx, entity.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword})

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Empty(t, out)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(out))
}

func createUpdatePasswordFailGetById(t *testing.T, ctx context.Context) (*application.ProfileApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		GetPassword(ctx, &gen.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword}).
		Return(&gen.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword, PasswordToCheck: oldPassword}, nil)
	profileService.EXPECT().
		UpdateUserPassword(ctx, &gen.UpdateFields{Id: id, User: &gen.User{Id: id, OldPassword: oldPassword, NewPassword: newPassword, PasswordToCheck: oldPassword}}).
		Return(&gen.Empty{}, nil)
	profileService.EXPECT().
		GetUserById(ctx, &gen.UserID{Id: id}).
		Return(&gen.PublicUser{}, errors.New("error"))

	securityService := mock.NewSecurityMock(ctrl)
	securityService.EXPECT().
		CheckPassword(oldPassword, oldPassword).
		Return(true)
	securityService.EXPECT().
		MakeShieldedPassword(newPassword).
		Return(newPassword, nil)

	return application.NewProfileApp(profileService, securityService), ctrl
}

func createContext() context.Context {
	return context.WithValue(context.Background(), entity.UserIdKey, id)
}

func createWatchlist() *gen.Currencies {
	return &gen.Currencies{
		Currencies: []*gen.Currency{
			{Base: curr, Title: base},
		},
	}
}
