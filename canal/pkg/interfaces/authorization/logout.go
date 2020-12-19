package authorization

import (
	"net/http"
	"server/canal/pkg/infrastructure/metric"
	"time"
)

func (a *Authentication) Logout(w http.ResponseWriter, r *http.Request) {
	defer metric.RecordTimeMetric(time.Now(), "Logout")

	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		metric.RecordHitMetric(http.StatusUnauthorized, r.URL.Path)
		return
	}

	desc, newCookie, err := a.authLogic.Logout(r.Context(), cookie)
	if err != nil {
		a.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		metric.RecordHitMetric(desc.Status, r.URL.Path)
		return
	}

	http.SetCookie(w, &newCookie)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	metric.RecordHitMetric(http.StatusOK, r.URL.Path)
}
