package application

import (
	"context"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"
	"net/http"
	"server/canal/pkg/domain/entity"
	profile "server/profile/pkg/profile/gen"
)

type ProfileApp struct {
	profile   profile.ProfileServiceClient
	sanitizer bluemonday.Policy
}

func (pfl *ProfileApp) UpdateAvatar(ctx context.Context, userEntity entity.User) (entity.Description, entity.PublicUser) {
	if !pfl.validate("Avatar", userEntity) {
		return entity.Description{
			Status:   http.StatusBadRequest,
			Function: "UpdateAvatar",
			Action:   "validate",
			Err:      errors.New("validation"),
		}, entity.PublicUser{}
	}

	userPfl := userEntity.ConvertToGRPC()
	if _, err := pfl.profile.UpdateUserAvatar(ctx, &profile.UpdateFields{Id: userPfl.Id, User: userPfl}); err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "UpdateAvatar",
			Action:   "UpdateUserAvatar",
			Err:      err,
		}, entity.PublicUser{}
	}

	public, err := pfl.profile.GetUserById(ctx, &profile.UserID{Id: userPfl.Id})
	if err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "UpdateAvatar",
			Action:   "GetUserById",
			Err:      err,
		}, entity.PublicUser{}
	}

	pfl.sanitizer.SanitizeBytes([]byte(public.Avatar))

	return entity.Description{}, entity.ConvertToPublic(public)
}

func (pfl *ProfileApp) validate(action string, user entity.User) bool {
	switch action {
	case "Avatar":
		if govalidator.IsNull(user.Avatar) {
			logrus.WithFields(logrus.Fields{
				"status":   http.StatusBadRequest,
				"function": "UpdateUserAvatar",
				"action":   "validation",
			}).Error("No user avatar from json")
			return false
		}
	case "Passwords":
		if govalidator.IsNull(user.OldPassword) || govalidator.IsNull(user.NewPassword) {
			logrus.WithFields(logrus.Fields{
				"status":      http.StatusBadRequest,
				"function":    "UpdateUserPassword",
				"oldPassword": user.OldPassword,
				"newPassword": user.NewPassword,
			}).Error("No user passwords from json")
			return false
		}
	}

	return true
}
