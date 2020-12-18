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
		errs.NonFieldError = append(errs.NonFieldError, "Неправильный email или пароль")
		return entity.Description{
			Status:    http.StatusBadRequest,
			Function:  "Login",
			Action:    "checkExistence",
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
	var userId = &authorization.UserID{Id: public.Id}
	if sess, err = authApp.auth.Create(ctx, userId); err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "Login",
			Action:   "Create",
		}, entity.PublicUser{}, http.Cookie{}, err
	}
	cookie.Value = sess.SessionId

	return entity.Description{}, entity.ConvertToPublic(public), cookie, nil
}

func (authApp *AuthApp) Logout(ctx context.Context, session authorization.Session) (entity.Description, error) {

}

func (authApp *AuthApp) Authenticate(ctx context.Context, userId authorization.UserID) (entity.Description, error) {

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
