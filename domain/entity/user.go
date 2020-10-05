package entity

type User struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PublicUser struct {
	Email string `json:"email"`
}
