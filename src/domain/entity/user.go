package entity

import (
	"github.com/asaskevich/govalidator"
	"server/src/domain/entity/jsonRealisation"
)

type User struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"required"`
	Avatar   string `json:"avatar"`
}

type PublicUser struct {
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

var Users []PublicUser

func (u *User) MakePublicUser() PublicUser {
	return PublicUser{
		Email:  u.Email,
		Avatar: u.Avatar,
	}
}

func (u *User) Validate(action string) (errors jsonRealisation.ErrorJSON) {
	switch action {
	case "signup":
		if !govalidator.IsEmail(u.Email) {
			errors.Email = append(errors.Email, "некорректный email")
			errors.NotEmpty = true
		}

		if u.Password == "" {
			errors.Password = append(errors.Email, "некорректный пароль")
			errors.NotEmpty = true
		}

		if govalidator.IsNull(u.Password) {
			errors.Password = append(errors.Email, "некорректный пароль")
			errors.NotEmpty = true
		}
	}

	return errors
}
