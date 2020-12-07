package logger

import "server/canal/pkg/domain/entity"

type Log interface {
	Info(desc entity.Description, err error)
	Warn(desc entity.Description, err error)
	Error(desc entity.Description, err error)
}
