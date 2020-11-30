package repository

type Utils interface {
	MakeShieldedPassword(string) (string, error)
	CheckPassword(string, string) bool
}
