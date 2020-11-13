package log

type LogHandler interface {
	Print(msg... interface{})
	Info(msg... interface{})
	Debug(msg... interface{})
}
