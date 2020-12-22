package profile

import (
	json "github.com/mailru/easyjson"
	"net/http"
	"server/canal/pkg/domain/entity"
)

func (p *Profile) createServerError(errs *entity.ErrorJSON, w http.ResponseWriter) int {
	res, err := json.Marshal(errs)
	if err != nil {
		code := http.StatusInternalServerError
		desc := entity.Description{Function: "createServerError", Action: "Marshal", Status: code}
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)
		return http.StatusInternalServerError
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	if _, err := w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "UpdateUserPassword", Action: "Write"}, err)
	}

	return http.StatusBadRequest
}
