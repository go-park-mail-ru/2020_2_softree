package entity

import (
	"github.com/asaskevich/govalidator"
)

type User struct {
	ID          uint64 `json:"id"`
	Email       string `json:"email" valid:"email"`
	Password    string `json:"password" valid:"required"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	Avatar      string `json:"avatar"`
}

type PublicUser struct {
	ID     uint64 `json:"id"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

var Users []PublicUser

func (u *User) MakePublicUser() PublicUser {
	return PublicUser{
		ID:     u.ID,
		Email:  u.Email,
		Avatar: u.Avatar,
	}
}

func (u *User) Validate() (errors ErrorJSON) {
	if !govalidator.IsEmail(u.Email) {
		errors.Email = append(errors.Email, "Некорректный email")
		errors.NotEmpty = true
	}

	if u.Password == "" {
		errors.Password = append(errors.Email, "Некорректный пароль")
		errors.NotEmpty = true
	}

	if govalidator.IsNull(u.Password) {
		errors.Password = append(errors.Email, "Некорректный пароль")
		errors.NotEmpty = true
	}

	if govalidator.HasWhitespace(u.Password) {
		errors.Password = append(errors.Email, "Некорректный пароль")
		errors.NotEmpty = true
	}

	return errors
}

func (u *User) ValidateUpdate() (errors ErrorJSON) {
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
