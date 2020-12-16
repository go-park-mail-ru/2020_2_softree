package entity

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"io"
	"io/ioutil"
	"net/http"
	profile "server/profile/pkg/profile/gen"
)

type User struct {
	Id              int64  `json:"id"`
	Email           string `json:"email" valid:"email"`
	Password        string `json:"password" valid:"required"`
	PasswordToCheck string `json:"password_to_check"`
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	Avatar          string `json:"avatar"`
}

type PublicUser struct {
	Id     int64  `json:"id"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

func GetUserFromBody(body io.ReadCloser) (User, Description, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return User{}, Description{Action: "ReadAll", Status: http.StatusInternalServerError}, err
	}
	defer body.Close()

	var user User
	err = user.UnmarshalJSON(data)
	if err != nil {
		return User{}, Description{Action: "Unmarshal", Status: http.StatusInternalServerError}, err
	}
	return user, Description{}, nil
}

func (user *User) Validate() error {
	check, err := govalidator.ValidateStruct(user)
	if err != nil {
		return err
	}
	if !check {
		return errors.New("validation")
	}
	return nil
}

func (user *User) ConvertToGRPC() *profile.User {
	return &profile.User{
		Id:              user.Id,
		Email:           user.Email,
		Password:        user.Password,
		PasswordToCheck: user.PasswordToCheck,
		OldPassword:     user.OldPassword,
		NewPassword:     user.NewPassword,
	}
}

func ConvertToUser(profileUser *profile.User) User {
	return User{
		Id:              profileUser.Id,
		Email:           profileUser.Email,
		Password:        profileUser.Password,
		PasswordToCheck: profileUser.PasswordToCheck,
		OldPassword:     profileUser.OldPassword,
		NewPassword:     profileUser.NewPassword,
		Avatar:          profileUser.Avatar,
	}
}

func ConvertToPublic(profileUser *profile.PublicUser) PublicUser {
	return PublicUser{
		Id:     profileUser.Id,
		Email:  profileUser.Email,
		Avatar: profileUser.Avatar,
	}
}
