package application

import (
	"server/canal/pkg/domain/entity"
)

func createErrorDescription(function, action string, code int) entity.Description {
	return entity.Description{
		Status: code,
		Function: function,
		Action: action,
	}
}
