package application

import (
	"context"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/microcosm-cc/bluemonday"
	"net/http"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/domain/repository"
	profile "server/profile/pkg/profile/gen"
)

type ProfileApp struct {
	profile   profile.ProfileServiceClient
	sanitizer bluemonday.Policy
	security  repository.Utils
}

func NewProfileApp(profile profile.ProfileServiceClient, security repository.Utils) *ProfileApp {
	return &ProfileApp{profile: profile, security: security, sanitizer: *bluemonday.UGCPolicy()}
}

func (pfl *ProfileApp) UpdateAvatar(ctx context.Context, userEntity entity.User) (entity.Description, entity.PublicUser, error) {
	if err := pfl.validate("Avatar", userEntity); err != nil {
		return entity.Description{
			Status:   http.StatusBadRequest,
			Function: "UpdateAvatar",
			Action:   "validate",
		}, entity.PublicUser{}, err
	}

	userPfl := userEntity.ConvertToGRPC()
	if _, err := pfl.profile.UpdateUserAvatar(ctx, &profile.UpdateFields{Id: userPfl.Id, User: userPfl}); err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "UpdateAvatar",
			Action:   "UpdateUserAvatar",
		}, entity.PublicUser{}, err
	}

	public, err := pfl.profile.GetUserById(ctx, &profile.UserID{Id: userPfl.Id})
	if err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "UpdateAvatar",
			Action:   "GetUserById",
		}, entity.PublicUser{}, err
	}

	pfl.sanitizer.SanitizeBytes([]byte(public.Avatar))

	return entity.Description{}, entity.ConvertToPublic(public), nil
}

func (pfl *ProfileApp) UpdatePassword(ctx context.Context, userEntity entity.User) (entity.Description, entity.PublicUser, error) {
	var err error
	if err := pfl.validate("Passwords", userEntity); err != nil {
		return entity.Description{
			Status:   http.StatusBadRequest,
			Function: "UpdatePassword",
			Action:   "validate",
		}, entity.PublicUser{}, err
	}

	pfl.sanitizer.Sanitize(userEntity.OldPassword)
	pfl.sanitizer.Sanitize(userEntity.NewPassword)

	if errs := pfl.validateUpdate(userEntity); errs.NotEmpty {
		return entity.Description{
			ErrorJSON: errs,
		}, entity.PublicUser{}, nil
	}

	user := userEntity.ConvertToGRPC()
	if user, err = pfl.profile.GetPassword(ctx, user); err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "UpdatePassword",
			Action:   "GetPassword",
		}, entity.PublicUser{}, err
	}
	if !pfl.security.CheckPassword(user.PasswordToCheck, user.OldPassword) {
		var errs entity.ErrorJSON
		errs.Password = append(errs.Password, "введен неверно старый пароль")
		errs.NotEmpty = true
		return entity.Description{
			ErrorJSON: errs,
		}, entity.PublicUser{}, err
	}

	if user.NewPassword, err = pfl.security.MakeShieldedPassword(user.NewPassword); err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "UpdatePassword",
			Action:   "MakeShieldedPassword",
		}, entity.PublicUser{}, err
	}

	if _, err = pfl.profile.UpdateUserPassword(ctx, &profile.UpdateFields{Id: user.Id, User: user}); err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "UpdatePassword",
			Action:   "UpdateUserPassword",
		}, entity.PublicUser{}, err
	}

	var public *profile.PublicUser
	if public, err = pfl.profile.GetUserById(ctx, &profile.UserID{Id: user.Id}); err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "UpdatePassword",
			Action:   "GetUserById",
		}, entity.PublicUser{}, err
	}

	return entity.Description{}, entity.ConvertToPublic(public), nil
}

func (pfl *ProfileApp) ReceiveUser(ctx context.Context, id int64) (entity.Description, entity.PublicUser, error) {
	var err error
	var public *profile.PublicUser
	if public, err = pfl.profile.GetUserById(ctx, &profile.UserID{Id: id}); err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "ReceiveUser",
			Action:   "GetUserById",
		}, entity.PublicUser{}, err
	}

	return entity.Description{}, entity.ConvertToPublic(public), nil
}

func (pfl *ProfileApp) ReceiveWatchlist(ctx context.Context, id int64) (entity.Description, entity.Currencies, error) {
	var err error
	var currencies *profile.Currencies
	if currencies, err = pfl.profile.GetUserWatchlist(ctx, &profile.UserID{Id: id}); err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "ReceiveUser",
			Action:   "GetUserWatchlist",
		}, entity.Currencies{}, err
	}

	return entity.Description{}, entity.ConvertToSlice(currencies), nil
}

func (pfl *ProfileApp) validate(action string, user entity.User) error {
	switch action {
	case "Avatar":
		if govalidator.IsNull(user.Avatar) {
			return errors.New("no user avatar from json")
		}
	case "Passwords":
		if govalidator.IsNull(user.OldPassword) || govalidator.IsNull(user.NewPassword) {
			return errors.New("no user passwords from json")
		}
	}

	return nil
}

func (pfl *ProfileApp) validateUpdate(u entity.User) (errors entity.ErrorJSON) {
	if govalidator.HasWhitespace(u.NewPassword) {
		errors.Password = append(errors.Email, "Некорректный новый пароль")
		errors.NotEmpty = true
	}

	if govalidator.HasWhitespace(u.OldPassword) {
		errors.Password = append(errors.Email, "Некорректный старый пароль")
		errors.NotEmpty = true
	}

	return errors
}
