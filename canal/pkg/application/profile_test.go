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
	id     = int64(1)
	email  = "hound@psina.ru"
	avatar = "base64"
	base   = "RUB"
	curr   = "USD"
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

func TestUpdateAvatar_Fail(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := createUpdateAvatarFail(t, ctx)
	defer ctrl.Finish()

	desc, out, err := testAuth.ReceiveWatchlist(ctx, id)

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Empty(t, out)
	require.Equal(t, http.StatusInternalServerError, desc.Status)
	require.Equal(t, reflect.TypeOf(entity.Currencies{}), reflect.TypeOf(out))
}

func createUpdateAvatarFail(t *testing.T, ctx context.Context) (*application.ProfileApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	profileService := profile.NewProfileMock(ctrl)
	profileService.EXPECT().
		GetUserWatchlist(ctx, &gen.UserID{Id: id}).
		Return(&gen.Currencies{}, errors.New("error"))

	securityService := mock.NewSecurityMock(ctrl)

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
