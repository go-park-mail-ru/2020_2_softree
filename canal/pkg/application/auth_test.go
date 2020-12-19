package application_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"reflect"
	authMock "server/authorization/pkg/infrastructure/mock"
	authorization "server/authorization/pkg/session/gen"
	"server/canal/pkg/application"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/infrastructure/mock"
	profileMock "server/profile/pkg/infrastructure/mock"
	profile "server/profile/pkg/profile/gen"
	"testing"
	"time"
)

func Test_Login_Success(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := create_Login_Success(t, ctx)
	defer ctrl.Finish()

	desc, publicUser, cookie, err := testAuth.Login(ctx, entity.User{Email: email, Password: password})

	require.NoError(t, err)
	require.Empty(t, desc)
	require.NotEmpty(t, publicUser)
	require.NotEmpty(t, cookie)
	require.Equal(t, cookie.Value, sessionId)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(publicUser))
}

func create_Login_Success(t *testing.T, ctx context.Context) (*application.AuthApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockUser := profileMock.NewProfileMock(ctrl)

	mockAuth.EXPECT().
		Create(ctx, &authorization.UserID{Id: id}).
		Return(&authorization.Session{Id: id, SessionId: sessionId}, nil)
	mockUser.EXPECT().
		GetUserByLogin(ctx, &profile.User{Email: email, Password: password}).
		Return(createExpectedUser(), nil)
	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email, Password: password}).
		Return(&profile.Check{Existence: true}, nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewAuthApp(mockUser, mockAuth, securityService), ctrl
}

func Test_Login_Fail_UserNotExist(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := create_Login_Fail_UserNotExists(t, ctx)
	defer ctrl.Finish()

	desc, publicUser, cookie, err := testAuth.Login(ctx, entity.User{Email: email, Password: password})

	require.NoError(t, err)
	require.NotEmpty(t, desc)
	require.Equal(t, http.StatusBadRequest, desc.Status)
	require.Equal(t, "CheckExistence", desc.Action)
	require.NotEmpty(t, desc.ErrorJSON)
	require.Equal(t, "Неправильный email или пароль", desc.ErrorJSON.NonFieldError[0])
	require.Empty(t, publicUser)
	require.Empty(t, cookie)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(publicUser))
}

func create_Login_Fail_UserNotExists(t *testing.T, ctx context.Context) (*application.AuthApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockUser := profileMock.NewProfileMock(ctrl)

	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email, Password: password}).
		Return(&profile.Check{Existence: false}, nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewAuthApp(mockUser, mockAuth, securityService), ctrl
}

func Test_Login_Fail_CheckExistanceFail(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := create_Login_Fail_CheckExistanceFail(t, ctx)
	defer ctrl.Finish()

	desc, publicUser, cookie, err := testAuth.Login(ctx, entity.User{Email: email, Password: password})

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Equal(t, http.StatusInternalServerError, desc.Status)
	require.Equal(t, "CheckExistence", desc.Action)
	require.Empty(t, publicUser)
	require.Empty(t, cookie)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(publicUser))
}

func create_Login_Fail_CheckExistanceFail(t *testing.T, ctx context.Context) (*application.AuthApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockUser := profileMock.NewProfileMock(ctrl)

	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email, Password: password}).
		Return(nil, errors.New("fail"))

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewAuthApp(mockUser, mockAuth, securityService), ctrl
}

func Test_Login_Fail_GetUserByLoginFail(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := create_Login_Fail_GetUserByLoginFail(t, ctx)
	defer ctrl.Finish()

	desc, publicUser, cookie, err := testAuth.Login(ctx, entity.User{Email: email, Password: password})

	require.NoError(t, err)
	require.NotEmpty(t, desc)
	require.Equal(t, http.StatusBadRequest, desc.Status)
	require.Equal(t, "GetUserByLogin", desc.Action)
	require.NotEmpty(t, desc.ErrorJSON)
	require.Empty(t, publicUser)
	require.Empty(t, cookie)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(publicUser))
}

func create_Login_Fail_GetUserByLoginFail(t *testing.T, ctx context.Context) (*application.AuthApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockUser := profileMock.NewProfileMock(ctrl)

	mockUser.EXPECT().
		GetUserByLogin(ctx, &profile.User{Email: email, Password: password}).
		Return(&profile.PublicUser{}, errors.New("fail"))
	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email, Password: password}).
		Return(&profile.Check{Existence: true}, nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewAuthApp(mockUser, mockAuth, securityService), ctrl
}

func Test_Logout_Success(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := create_Logout_Success(t, ctx)
	defer ctrl.Finish()

	cookie := http.Cookie{Value: sessionId}

	desc, cookie, err := testAuth.Logout(ctx, &cookie)

	require.NoError(t, err)
	require.Empty(t, desc)
	require.NotEmpty(t, cookie)
	require.Equal(t, time.Date(1973, 1, 1, 0, 0, 0, 0, time.UTC), cookie.Expires)
}

func create_Logout_Success(t *testing.T, ctx context.Context) (*application.AuthApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockUser := profileMock.NewProfileMock(ctrl)

	mockAuth.EXPECT().
		Delete(ctx, &authorization.SessionID{SessionId: sessionId}).
		Return(&authorization.Empty{}, nil)

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewAuthApp(mockUser, mockAuth, securityService), ctrl
}

func Test_Logout_Fail(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := create_Logout_Fail(t, ctx)
	defer ctrl.Finish()

	cookie := http.Cookie{Value: sessionId}

	desc, cookie, err := testAuth.Logout(ctx, &cookie)

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Equal(t, http.StatusInternalServerError, desc.Status)
	require.Equal(t, "Delete", desc.Action)
	require.Empty(t, cookie)
}

func create_Logout_Fail(t *testing.T, ctx context.Context) (*application.AuthApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockUser := profileMock.NewProfileMock(ctrl)

	mockAuth.EXPECT().
		Delete(ctx, &authorization.SessionID{SessionId: sessionId}).
		Return(&authorization.Empty{}, errors.New("fail"))

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewAuthApp(mockUser, mockAuth, securityService), ctrl
}

func Test_Signup_Success(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := create_Signup_Success(t, ctx)
	defer ctrl.Finish()

	desc, publicUser, cookie, err := testAuth.Signup(ctx, entity.User{Email: email, Password: password})

	require.NoError(t, err)
	require.Empty(t, desc)
	require.NotEmpty(t, cookie)
	require.Equal(t, email, publicUser.Email)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(publicUser))
}

func create_Signup_Success(t *testing.T, ctx context.Context) (*application.AuthApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockUser := profileMock.NewProfileMock(ctrl)
	securityService := mock.NewSecurityMock(ctrl)

	mockAuth.EXPECT().
		Create(ctx, &authorization.UserID{Id: id}).
		Return(&authorization.Session{Id: id, SessionId: sessionId}, nil)
	mockUser.EXPECT().
		PutPortfolio(ctx, &profile.PortfolioValue{Id: id, Value: amount}).
		Return(&profile.Empty{}, nil)
	mockUser.EXPECT().
		CreateInitialWallet(ctx, &profile.UserID{Id: id}).
		Return(&profile.Empty{}, nil)
	mockUser.EXPECT().
		SaveUser(ctx, &profile.User{Email: email, Password: password}).
		Return(createExpectedUser(), nil)
	securityService.EXPECT().
		MakeShieldedPassword(password).
		Return(password, nil)
	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email, Password: password}).
		Return(&profile.Check{Existence: false}, nil)

	return application.NewAuthApp(mockUser, mockAuth, securityService), ctrl
}

func Test_Signup_Fail_UserExists(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := create_Signup_Fail_UserExists(t, ctx)
	defer ctrl.Finish()

	desc, publicUser, cookie, err := testAuth.Signup(ctx, entity.User{Email: email, Password: password})

	require.NoError(t, err)
	require.NotEmpty(t, desc)
	require.Equal(t, http.StatusBadRequest, desc.Status)
	require.Equal(t, "CheckExistence", desc.Action)
	require.NotEmpty(t, desc.ErrorJSON)
	require.Equal(t, "Пользователь с таким email'ом уже существует", desc.ErrorJSON.NonFieldError[0])
	require.Empty(t, publicUser)
	require.Empty(t, cookie)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(publicUser))
}

func create_Signup_Fail_UserExists(t *testing.T, ctx context.Context) (*application.AuthApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockUser := profileMock.NewProfileMock(ctrl)
	securityService := mock.NewSecurityMock(ctrl)

	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email, Password: password}).
		Return(&profile.Check{Existence: true}, nil)

	return application.NewAuthApp(mockUser, mockAuth, securityService), ctrl
}

func Test_Signup_Fail_CheckExistenceFail(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := create_Signup_Fail_CheckExistenceFail(t, ctx)
	defer ctrl.Finish()

	desc, publicUser, cookie, err := testAuth.Signup(ctx, entity.User{Email: email, Password: password})

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Equal(t, http.StatusInternalServerError, desc.Status)
	require.Equal(t, "CheckExistence", desc.Action)
	require.Empty(t, publicUser)
	require.Empty(t, cookie)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(publicUser))
}

func create_Signup_Fail_CheckExistenceFail(t *testing.T, ctx context.Context) (*application.AuthApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockUser := profileMock.NewProfileMock(ctrl)

	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email, Password: password}).
		Return(nil, errors.New("fail"))

	securityService := mock.NewSecurityMock(ctrl)

	return application.NewAuthApp(mockUser, mockAuth, securityService), ctrl
}

func Test_Signup_Fail_MakeShieldedPasswordFail(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := create_Signup_Fail_MakeShieldedPasswordFail(t, ctx)
	defer ctrl.Finish()

	desc, publicUser, cookie, err := testAuth.Signup(ctx, entity.User{Email: email, Password: password})

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Equal(t, http.StatusInternalServerError, desc.Status)
	require.Equal(t, "MakeShieldedPassword", desc.Action)
	require.Empty(t, publicUser)
	require.Empty(t, cookie)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(publicUser))
}

func create_Signup_Fail_MakeShieldedPasswordFail(t *testing.T, ctx context.Context) (*application.AuthApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockUser := profileMock.NewProfileMock(ctrl)
	securityService := mock.NewSecurityMock(ctrl)

	securityService.EXPECT().
		MakeShieldedPassword(password).
		Return("", errors.New("fail"))
	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email, Password: password}).
		Return(&profile.Check{Existence: false}, nil)

	return application.NewAuthApp(mockUser, mockAuth, securityService), ctrl
}

func Test_Signup_Fail_SaveUserFail(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := create_Signup_Fail_SaveUserFail(t, ctx)
	defer ctrl.Finish()

	desc, publicUser, cookie, err := testAuth.Signup(ctx, entity.User{Email: email, Password: password})

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Equal(t, http.StatusInternalServerError, desc.Status)
	require.Equal(t, "SaveUser", desc.Action)
	require.Empty(t, publicUser)
	require.Empty(t, cookie)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(publicUser))
}

func create_Signup_Fail_SaveUserFail(t *testing.T, ctx context.Context) (*application.AuthApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockUser := profileMock.NewProfileMock(ctrl)
	securityService := mock.NewSecurityMock(ctrl)

	mockUser.EXPECT().
		SaveUser(ctx, &profile.User{Email: email, Password: password}).
		Return(&profile.PublicUser{}, errors.New("fail"))
	securityService.EXPECT().
		MakeShieldedPassword(password).
		Return(password, nil)
	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email, Password: password}).
		Return(&profile.Check{Existence: false}, nil)

	return application.NewAuthApp(mockUser, mockAuth, securityService), ctrl
}

func Test_Signup_Fail_CreateInitialWalletFail(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := create_Signup_Fail_CreateInitialWalletFail(t, ctx)
	defer ctrl.Finish()

	desc, publicUser, cookie, err := testAuth.Signup(ctx, entity.User{Email: email, Password: password})

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Equal(t, http.StatusInternalServerError, desc.Status)
	require.Equal(t, "CreateInitialWallet", desc.Action)
	require.Empty(t, publicUser)
	require.Empty(t, cookie)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(publicUser))
}

func create_Signup_Fail_CreateInitialWalletFail(t *testing.T, ctx context.Context) (*application.AuthApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockUser := profileMock.NewProfileMock(ctrl)
	securityService := mock.NewSecurityMock(ctrl)

	mockUser.EXPECT().
		CreateInitialWallet(ctx, &profile.UserID{Id: id}).
		Return(&profile.Empty{}, errors.New("fail"))
	mockUser.EXPECT().
		SaveUser(ctx, &profile.User{Email: email, Password: password}).
		Return(createExpectedUser(), nil)
	securityService.EXPECT().
		MakeShieldedPassword(password).
		Return(password, nil)
	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email, Password: password}).
		Return(&profile.Check{Existence: false}, nil)

	return application.NewAuthApp(mockUser, mockAuth, securityService), ctrl
}

func Test_Signup_Fail_PutPortfolioFail(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := create_Signup_Fail_PutPortfolioFail(t, ctx)
	defer ctrl.Finish()

	desc, publicUser, cookie, err := testAuth.Signup(ctx, entity.User{Email: email, Password: password})

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Equal(t, http.StatusInternalServerError, desc.Status)
	require.Equal(t, "PutPortfolio", desc.Action)
	require.Empty(t, publicUser)
	require.Empty(t, cookie)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(publicUser))
}

func create_Signup_Fail_PutPortfolioFail(t *testing.T, ctx context.Context) (*application.AuthApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockUser := profileMock.NewProfileMock(ctrl)
	securityService := mock.NewSecurityMock(ctrl)

	mockUser.EXPECT().
		PutPortfolio(ctx, &profile.PortfolioValue{Id: id, Value: amount}).
		Return(&profile.Empty{}, errors.New("fail"))
	mockUser.EXPECT().
		CreateInitialWallet(ctx, &profile.UserID{Id: id}).
		Return(&profile.Empty{}, nil)
	mockUser.EXPECT().
		SaveUser(ctx, &profile.User{Email: email, Password: password}).
		Return(createExpectedUser(), nil)
	securityService.EXPECT().
		MakeShieldedPassword(password).
		Return(password, nil)
	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email, Password: password}).
		Return(&profile.Check{Existence: false}, nil)

	return application.NewAuthApp(mockUser, mockAuth, securityService), ctrl
}

func Test_Signup_Fail_CreateFail(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := create_Signup_Fail_CreateFail(t, ctx)
	defer ctrl.Finish()

	desc, publicUser, cookie, err := testAuth.Signup(ctx, entity.User{Email: email, Password: password})

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Equal(t, http.StatusInternalServerError, desc.Status)
	require.Equal(t, "Create", desc.Action)
	require.Empty(t, publicUser)
	require.Empty(t, cookie)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(publicUser))
}

func create_Signup_Fail_CreateFail(t *testing.T, ctx context.Context) (*application.AuthApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockUser := profileMock.NewProfileMock(ctrl)
	securityService := mock.NewSecurityMock(ctrl)

	mockAuth.EXPECT().
		Create(ctx, &authorization.UserID{Id: id}).
		Return(&authorization.Session{}, errors.New("fail"))
	mockUser.EXPECT().
		PutPortfolio(ctx, &profile.PortfolioValue{Id: id, Value: amount}).
		Return(&profile.Empty{}, nil)
	mockUser.EXPECT().
		CreateInitialWallet(ctx, &profile.UserID{Id: id}).
		Return(&profile.Empty{}, nil)
	mockUser.EXPECT().
		SaveUser(ctx, &profile.User{Email: email, Password: password}).
		Return(createExpectedUser(), nil)
	securityService.EXPECT().
		MakeShieldedPassword(password).
		Return(password, nil)
	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email, Password: password}).
		Return(&profile.Check{Existence: false}, nil)

	return application.NewAuthApp(mockUser, mockAuth, securityService), ctrl
}

func Test_Authenticate_Success(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := create_Authenticate_Success(t, ctx)
	defer ctrl.Finish()

	desc, publicUser, err := testAuth.Authenticate(ctx, id)

	require.NoError(t, err)
	require.Empty(t, desc)
	require.NotEmpty(t, publicUser)
	require.Equal(t, email, publicUser.Email)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(publicUser))
}

func create_Authenticate_Success(t *testing.T, ctx context.Context) (*application.AuthApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockUser := profileMock.NewProfileMock(ctrl)
	securityService := mock.NewSecurityMock(ctrl)

	mockUser.EXPECT().
		GetUserById(ctx, &profile.UserID{Id: id}).
		Return(createExpectedUser(), nil)

	return application.NewAuthApp(mockUser, mockAuth, securityService), ctrl
}

func Test_Authenticate_Fail(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := create_Authenticate_Fail(t, ctx)
	defer ctrl.Finish()

	desc, publicUser, err := testAuth.Authenticate(ctx, id)

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Equal(t, http.StatusBadRequest, desc.Status)
	require.Equal(t, "GetUserById", desc.Action)
	require.Empty(t, publicUser)
	require.Equal(t, reflect.TypeOf(entity.PublicUser{}), reflect.TypeOf(publicUser))
}

func create_Authenticate_Fail(t *testing.T, ctx context.Context) (*application.AuthApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockUser := profileMock.NewProfileMock(ctrl)
	securityService := mock.NewSecurityMock(ctrl)

	mockUser.EXPECT().
		GetUserById(ctx, &profile.UserID{Id: id}).
		Return(&profile.PublicUser{}, errors.New("fail"))

	return application.NewAuthApp(mockUser, mockAuth, securityService), ctrl
}

func Test_Auth_Success(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := create_Auth_Success(t, ctx)
	defer ctrl.Finish()

	cookie := http.Cookie{Value: sessionId}

	desc, userId, err := testAuth.Auth(ctx, &cookie)

	require.NoError(t, err)
	require.Empty(t, desc)
	require.Equal(t, id, userId)
}

func create_Auth_Success(t *testing.T, ctx context.Context) (*application.AuthApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockUser := profileMock.NewProfileMock(ctrl)
	securityService := mock.NewSecurityMock(ctrl)

	mockAuth.EXPECT().
		Check(ctx, &authorization.SessionID{SessionId: sessionId}).
		Return(&authorization.UserID{Id: id}, nil)

	return application.NewAuthApp(mockUser, mockAuth, securityService), ctrl
}

func Test_Auth_Fail(t *testing.T) {
	ctx := createContext()
	testAuth, ctrl := create_Auth_Fail(t, ctx)
	defer ctrl.Finish()

	cookie := http.Cookie{Value: sessionId}

	desc, userId, err := testAuth.Auth(ctx, &cookie)

	require.Error(t, err)
	require.NotEmpty(t, desc)
	require.Equal(t, http.StatusBadRequest, desc.Status)
	require.Equal(t, "Check", desc.Action)
	require.Equal(t, int64(0), userId)
}

func create_Auth_Fail(t *testing.T, ctx context.Context) (*application.AuthApp, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockUser := profileMock.NewProfileMock(ctrl)
	securityService := mock.NewSecurityMock(ctrl)

	mockAuth.EXPECT().
		Check(ctx, &authorization.SessionID{SessionId: sessionId}).
		Return(&authorization.UserID{}, errors.New("fail"))

	return application.NewAuthApp(mockUser, mockAuth, securityService), ctrl
}

func createExpectedUser() *profile.PublicUser {
	return &profile.PublicUser{Id: id, Email: email, Avatar: avatar}
}
