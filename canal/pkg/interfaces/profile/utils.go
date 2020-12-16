package profile

import (
	json "github.com/mailru/easyjson"
	"net/http"
	"server/canal/pkg/domain/entity"
	"strconv"
)

func (p *Profile) createErrorJSON(err error) (errs entity.ErrorJSON) {
	if err.Error() == "wrong old password" {
		errs.NotEmpty = true
		errs.NonFieldError = append(errs.NonFieldError, "Введен неверно старый пароль")
	}

	if err.Error() == "not enough payment" {
		errs.NotEmpty = true
		errs.NonFieldError = append(
			errs.NonFieldError,
			"В вашем кошельке недостаточно средств для совершения данной транзакции",
		)
	}

	return
}

func (p *Profile) createServerError(errs *entity.ErrorJSON, w http.ResponseWriter) {
	res, err := json.Marshal(errs)
	if err != nil {
		code := http.StatusInternalServerError
		desc := entity.Description{Function: "createServerError", Action: "Marshal", Status: code}
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	p.recordHitMetric(http.StatusBadRequest)
	if _, err := w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "UpdateUserPassword", Action: "Write"}, err)
	}
}

func (p *Profile) recordHitMetric(code int) {
	p.Hits.WithLabelValues(strconv.Itoa(code)).Inc()
}


