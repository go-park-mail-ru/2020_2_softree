package authorization

import (
	json "github.com/mailru/easyjson"
	"net/http"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/infrastructure/metric"
	"time"
)

func (a *Authentication) Login(w http.ResponseWriter, r *http.Request) {
	defer metric.RecordTimeMetric(time.Now(), "Login")

	user, desc, err := entity.GetUserFromBody(r.Body)
	if err != nil {
		desc.Function = "Login"
		a.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}

	desc, public, cookie, err := a.authLogic.Login(r.Context(), user)
	if err != nil {
		a.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		metric.RecordHitMetric(desc.Status, r.URL.Path)
		return
	}

	if desc.ErrorJSON.NotEmpty {
		a.handleErrorJSON(desc, w)

		metric.RecordHitMetric(desc.Status, r.URL.Path)
		return
	}
	http.SetCookie(w, &cookie)

	res, err := json.Marshal(public)
	if err != nil {
		desc = entity.Description{Function: "Login",
			Action: "Marshal",
			Status: http.StatusInternalServerError}
		a.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	metric.RecordHitMetric(http.StatusOK, r.URL.Path)
	if _, err := w.Write(res); err != nil {
		a.logger.Error(entity.Description{Function: "Login", Action: "Write"}, err)
	}
}
