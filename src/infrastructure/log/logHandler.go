package log

type LogHandler interface {
	Info(msg... interface{})
	Debug(msg... interface{})
}
