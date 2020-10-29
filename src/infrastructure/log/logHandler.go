package log

type LogHandler interface {
	Print(err interface{})
}
