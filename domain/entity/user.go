package entity

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
