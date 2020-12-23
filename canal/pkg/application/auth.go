package application

import (
	"context"
	"github.com/asaskevich/govalidator"
	"github.com/microcosm-cc/bluemonday"
	"net/http"
	authorization "server/authorization/pkg/session/gen"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/domain/repository"
	utils "server/canal/pkg/interfaces/authorization"
	profile "server/profile/pkg/profile/gen"
	"time"
)

type AuthApp struct {
	profile   profile.ProfileServiceClient
	auth      authorization.AuthorizationServiceClient
	sanitizer bluemonday.Policy
	security  repository.Utils
}

func NewAuthApp(profile profile.ProfileServiceClient, auth authorization.AuthorizationServiceClient, security repository.Utils) *AuthApp {
	return &AuthApp{profile: profile, auth: auth, sanitizer: *bluemonday.UGCPolicy(), security: security}
}

func (authApp *AuthApp) Login(ctx context.Context, user entity.User) (entity.Description, entity.PublicUser, http.Cookie, error) {
	authApp.sanitizer.Sanitize(user.Email)
	authApp.sanitizer.Sanitize(user.Password)

	var errs entity.ErrorJSON
	var desc entity.Description
	if errs = authApp.validate(user); errs.NotEmpty {
		desc = createErrorDescription("Login", "validate", http.StatusBadRequest)
		desc.ErrorJSON = errs
		return desc, entity.PublicUser{}, http.Cookie{}, nil
	}

	userGRPC := user.ConvertToGRPC()

	check, err := authApp.profile.CheckExistence(ctx, userGRPC)
	if err != nil {
		return createErrorDescription("Login", "CheckExistence", http.StatusInternalServerError),
			entity.PublicUser{}, http.Cookie{}, err
	}
	if !check.Existence {
		errs.NotEmpty = true
		errs.NonFieldError = append(errs.NonFieldError, "Неправильный email или пароль")

		desc = createErrorDescription("Login", "CheckExistence", http.StatusBadRequest)
		desc.ErrorJSON = errs

		return desc, entity.PublicUser{}, http.Cookie{}, nil
	}

	public, err := authApp.profile.GetUserByLogin(ctx, userGRPC)
	if err != nil {
		if errs = authApp.checkGetUserByLoginErrors(err); errs.NotEmpty {
			desc = createErrorDescription("Login", "GetUserByLogin", http.StatusBadRequest)
			desc.ErrorJSON = errs
			return desc, entity.PublicUser{}, http.Cookie{}, nil
		}
	}

	cookie := utils.CreateCookie()
	var sess *authorization.Session
	if sess, err = authApp.auth.Create(ctx, &authorization.UserID{Id: public.Id}); err != nil {
		return createErrorDescription("Login", "Create", http.StatusInternalServerError),
			entity.PublicUser{}, http.Cookie{}, err
	}
	cookie.Value = sess.SessionId

	return entity.Description{}, entity.ConvertToPublic(public), cookie, nil
}

func (authApp *AuthApp) Logout(ctx context.Context, cookie *http.Cookie) (entity.Description, http.Cookie, error) {
	if _, err := authApp.auth.Delete(ctx, &authorization.SessionID{SessionId: cookie.Value}); err != nil {
		return createErrorDescription("Logout", "Delete", http.StatusInternalServerError),
			http.Cookie{}, err
	}

	newCookie := utils.CreateCookie()
	newCookie.Expires = time.Date(1973, 1, 1, 0, 0, 0, 0, time.UTC)
	newCookie.Value = ""

	return entity.Description{}, newCookie, nil
}

func (authApp *AuthApp) Signup(ctx context.Context, user entity.User) (entity.Description, entity.PublicUser, http.Cookie, error) {
	authApp.sanitizer.Sanitize(user.Email)
	authApp.sanitizer.Sanitize(user.Password)

	var errs entity.ErrorJSON
	var desc entity.Description
	if errs = authApp.validate(user); errs.NotEmpty {
		desc = createErrorDescription("Signup", "validate", http.StatusBadRequest)
		desc.ErrorJSON = errs
		return desc, entity.PublicUser{}, http.Cookie{}, nil
	}

	userGRPC := user.ConvertToGRPC()

	check, err := authApp.profile.CheckExistence(ctx, userGRPC)
	if err != nil {
		return createErrorDescription("Signup", "CheckExistence", http.StatusInternalServerError),
			entity.PublicUser{}, http.Cookie{}, err
	}
	if check.Existence {
		errs.NotEmpty = true
		errs.NonFieldError = append(errs.NonFieldError, "Пользователь с таким email'ом уже существует")

		desc = createErrorDescription("Signup", "CheckExistence", http.StatusBadRequest)
		desc.ErrorJSON = errs

		return desc, entity.PublicUser{}, http.Cookie{}, nil
	}

	if userGRPC.Password, err = authApp.security.MakeShieldedPassword(userGRPC.Password); err != nil {
		return createErrorDescription("Signup", "MakeShieldedPassword", http.StatusInternalServerError),
			entity.PublicUser{}, http.Cookie{}, err
	}

	public, err := authApp.profile.SaveUser(ctx, userGRPC)
	if err != nil {
		return createErrorDescription("Signup", "SaveUser", http.StatusInternalServerError),
			entity.PublicUser{}, http.Cookie{}, err
	}

	if _, err = authApp.profile.CreateInitialWallet(ctx, &profile.UserID{Id: public.Id}); err != nil {
		return createErrorDescription("Signup", "CreateInitialWallet", http.StatusInternalServerError),
			entity.PublicUser{}, http.Cookie{}, err
	}

	if _, err = authApp.profile.PutPortfolio(ctx, &profile.PortfolioValue{Id: public.Id, Value: 1000}); err != nil {
		return createErrorDescription("Signup", "PutPortfolio", http.StatusInternalServerError),
			entity.PublicUser{}, http.Cookie{}, err
	}

	cookie := utils.CreateCookie()
	var sess *authorization.Session
	if sess, err = authApp.auth.Create(ctx, &authorization.UserID{Id: public.Id}); err != nil {
		return createErrorDescription("Signup", "Create", http.StatusInternalServerError),
			entity.PublicUser{}, http.Cookie{}, err
	}
	cookie.Value = sess.SessionId

	return entity.Description{}, entity.ConvertToPublic(public), cookie, nil
}

func (authApp *AuthApp) Authenticate(ctx context.Context, userId int64) (entity.Description, entity.PublicUser, error) {
	user, err := authApp.profile.GetUserById(ctx, &profile.UserID{Id: userId})
	if err != nil {
		return createErrorDescription("Authenticate", "GetUserById", http.StatusInternalServerError),
			entity.PublicUser{}, err
	}

	return entity.Description{}, entity.ConvertToPublic(user), nil
}

func (authApp *AuthApp) Auth(ctx context.Context, cookie *http.Cookie) (entity.Description, int64, error) {
	id, err := authApp.auth.Check(ctx, &authorization.SessionID{SessionId: cookie.Value})
	if err != nil {
		return createErrorDescription("Auth", "Check", http.StatusBadRequest), 0, err
	}

	return entity.Description{}, id.Id, nil
}

func (authApp *AuthApp) validate(user entity.User) (errs entity.ErrorJSON) {
	if !govalidator.IsEmail(user.Email) {
		errs.Email = append(errs.Email, "Некорректный email или пароль")
		errs.NotEmpty = true
	}

	if user.Password == "" {
		errs.Password = append(errs.Email, "Некорректный email или пароль")
		errs.NotEmpty = true
	}

	if govalidator.IsNull(user.Password) {
		errs.Password = append(errs.Email, "Некорректный email или пароль")
		errs.NotEmpty = true
	}

	if govalidator.HasWhitespace(user.Password) {
		errs.Password = append(errs.Email, "Некорректный email или пароль")
		errs.NotEmpty = true
	}

	return
}

func (authApp *AuthApp) checkGetUserByLoginErrors(err error) (errs entity.ErrorJSON) {
	if err != nil {
		errs.NotEmpty = true
		errs.NonFieldError = append(errs.NonFieldError, "")
	}

	if err.Error() == "wrong password" {
		errs.NotEmpty = true
		errs.NonFieldError = append(errs.NonFieldError, "Неправильный email или пароль")
	}

	return
}
