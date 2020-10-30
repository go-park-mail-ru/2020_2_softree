package auth

type AuthHandler interface {
	AuthCreator
	AuthChecker
	AuthEraser
}

type AuthCreator interface {
	CreateAuth(uint64, string) error
}

type AuthChecker interface {
	CheckAuth(string) (uint64, error)
}

type AuthEraser interface {
	DeleteAuth(*AccessDetails) error
}

type AccessDetails struct {
	Value string `json:"session_id"`
}
