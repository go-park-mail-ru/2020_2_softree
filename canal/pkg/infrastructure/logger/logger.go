package logger

import "server/canal/pkg/domain/entity"

type Log interface {
	Info(desc entity.Description)
	Warn(desc entity.Description)
	Error(desc entity.Description)
}
