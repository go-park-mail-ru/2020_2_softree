package auth

type AuthInterface interface {
	CreateAuth(uint64, *CookieDetails) error
	CheckAuth(string) (uint64, error)
	DeleteAuth(*AccessDetails) error
}

type CookieDetails struct {
}

type AccessDetails struct {
	SessionId string
	UserId    uint64
}
