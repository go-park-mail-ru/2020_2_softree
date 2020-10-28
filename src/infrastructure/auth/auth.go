package auth

type AuthHandler interface {
	CreateAuth(uint64, string) error
	CheckAuth(string) (uint64, error)
	DeleteAuth(*AccessDetails) error
}

type AccessDetails struct {
	Value string `json:"session_id"`
}
