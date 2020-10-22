package auth

type AuthInterface interface {
	CreateAuth(uint64, string) error
	CheckAuth(string) (uint64, error)
	DeleteAuth(*AccessDetails) error
}

type AccessDetails struct {
	SessionId uint64 `json:"session_id"`
}
