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
	if errs = authApp.validate(user); errs.NotEmpty {
		return entity.Description{
			Status:    http.StatusBadRequest,
			Function:  "Login",
			Action:    "validate",
			ErrorJSON: errs,
		}, entity.PublicUser{}, http.Cookie{}, nil
	}

	userGRPC := user.ConvertToGRPC()

	check, err := authApp.profile.CheckExistence(ctx, userGRPC)
	if err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "Login",
			Action:   "CheckExistence",
		}, entity.PublicUser{}, http.Cookie{}, err
	}
	if !check.Existence {
		errs.NotEmpty = true
		errs.NonFieldError = append(errs.NonFieldError, "Неправильный email или пароль")
		return entity.Description{
			Status:    http.StatusBadRequest,
			Function:  "Login",
			Action:    "CheckExistence",
			ErrorJSON: errs,
		}, entity.PublicUser{}, http.Cookie{}, nil
	}

	public, err := authApp.profile.GetUserByLogin(ctx, userGRPC)
	if err != nil {
		if errs = authApp.checkGetUserByLoginErrors(err); errs.NotEmpty {
			return entity.Description{
				Status:    http.StatusBadRequest,
				Function:  "Login",
				Action:    "GetUserByLogin",
				ErrorJSON: errs,
			}, entity.PublicUser{}, http.Cookie{}, nil
		}
	}

	cookie := utils.CreateCookie()
	var sess *authorization.Session
	if sess, err = authApp.auth.Create(ctx, &authorization.UserID{Id: public.Id}); err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "Login",
			Action:   "Create",
		}, entity.PublicUser{}, http.Cookie{}, err
	}
	cookie.Value = sess.SessionId

	return entity.Description{}, entity.ConvertToPublic(public), cookie, nil
}

func (authApp *AuthApp) Logout(ctx context.Context, cookie *http.Cookie) (entity.Description, http.Cookie, error) {
	if _, err := authApp.auth.Delete(ctx, &authorization.SessionID{SessionId: cookie.Value}); err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "Logout",
			Action:   "Delete",
		}, http.Cookie{}, err
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
	if errs = authApp.validate(user); errs.NotEmpty {
		return entity.Description{
			Status:    http.StatusBadRequest,
			Function:  "Signup",
			Action:    "validate",
			ErrorJSON: errs,
		}, entity.PublicUser{}, http.Cookie{}, nil
	}

	userGRPC := user.ConvertToGRPC()

	check, err := authApp.profile.CheckExistence(ctx, userGRPC)
	if err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "Signup",
			Action:   "CheckExistence",
		}, entity.PublicUser{}, http.Cookie{}, err
	}
	if check.Existence {
		errs.NotEmpty = true
		errs.NonFieldError = append(errs.NonFieldError, "Пользователь с таким email'ом уже существует")
		return entity.Description{
			Status:    http.StatusBadRequest,
			Function:  "Signup",
			Action:    "CheckExistence",
			ErrorJSON: errs,
		}, entity.PublicUser{}, http.Cookie{}, nil
	}

	if user.Password, err = authApp.security.MakeShieldedPassword(user.Password); err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "Signup",
			Action:   "MakeShieldedPassword",
		}, entity.PublicUser{}, http.Cookie{}, err
	}

	public, err := authApp.profile.SaveUser(ctx, userGRPC)
	if err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "Signup",
			Action:   "SaveUser",
		}, entity.PublicUser{}, http.Cookie{}, err
	}

	if _, err = authApp.profile.CreateInitialWallet(ctx, &profile.UserID{Id: public.Id}); err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "Signup",
			Action:   "CreateInitialWallet",
		}, entity.PublicUser{}, http.Cookie{}, err
	}

	if _, err = authApp.profile.PutPortfolio(ctx, &profile.PortfolioValue{Id: public.Id, Value: 1000}); err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "Signup",
			Action:   "PutPortfolio",
		}, entity.PublicUser{}, http.Cookie{}, err
	}

	cookie := utils.CreateCookie()
	var sess *authorization.Session
	if sess, err = authApp.auth.Create(ctx, &authorization.UserID{Id: public.Id}); err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "Signup",
			Action:   "Create",
		}, entity.PublicUser{}, http.Cookie{}, err
	}
	cookie.Value = sess.SessionId

	return entity.Description{}, entity.ConvertToPublic(public), cookie, nil
}

func (authApp *AuthApp) Authenticate(ctx context.Context, userId int64) (entity.Description, entity.PublicUser, error) {
	user, err := authApp.profile.GetUserById(ctx, &profile.UserID{Id: userId})
	if err != nil {
		return entity.Description{
			Status:   http.StatusBadRequest,
			Function: "Authenticate",
			Action:   "GetUserById",
		}, entity.PublicUser{}, err
	}

	return entity.Description{}, entity.ConvertToPublic(user), nil
}

func (authApp *AuthApp) Auth(ctx context.Context, cookie *http.Cookie) (entity.Description, int64, error) {
	id, err := authApp.auth.Check(ctx, &authorization.SessionID{SessionId: cookie.Value})
	if err != nil {
		return entity.Description{
			Status:   http.StatusBadRequest,
			Function: "Auth",
			Action:   "Check",
		}, 0, err
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
