package entity

import "server/src/domain/entity/jsonRealisation"

type User struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

type PublicUser struct {
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

var Users []PublicUser

func (u *User) PublicUsers() PublicUser {
	return PublicUser{
		Email:  u.Email,
		Avatar: u.Avatar,
	}
}

func (u *User) Validate(action string) jsonRealisation.ErrorJSON {
	// some user validation like email, password, password difference
	// action like login, auth, signup and others
	// returns errorJSON
	// errorJSON will be converted to json from calling func
	return jsonRealisation.ErrorJSON{}
}
