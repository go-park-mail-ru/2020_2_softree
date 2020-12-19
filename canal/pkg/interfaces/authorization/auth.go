package authorization

import (
	"context"
	json "github.com/mailru/easyjson"
	"net/http"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/infrastructure/metric"
	"time"
)

func (a *Authentication) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			metric.RecordHitMetric(http.StatusUnauthorized, r.URL.Path)
			return
		}
		if err != nil {
			desc := entity.Description{
				Function: "AuthMiddleware",
				Action:   "r.Cookie()",
				Status:   http.StatusBadRequest}
			a.logger.Error(desc, err)
			w.WriteHeader(http.StatusBadRequest)

			metric.RecordHitMetric(http.StatusBadRequest, r.URL.Path)
			return
		}

		if cookie == nil {
			desc := entity.Description{
				Function: "AuthMiddleware",
				Action:   "if cookie == nil",
				Status:   http.StatusBadRequest}
			a.logger.Error(desc, err)
			w.WriteHeader(http.StatusBadRequest)

			metric.RecordHitMetric(http.StatusBadRequest, r.URL.Path)
			return
		}

		desc, id, err := a.authLogic.Auth(r.Context(), cookie)
		if err != nil {
			a.logger.Error(desc, err)
			w.WriteHeader(desc.Status)

			metric.RecordHitMetric(desc.Status, r.URL.Path)
			return
		}

		ctx := context.WithValue(r.Context(), entity.UserIdKey, id)
		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	}
}

func (a *Authentication) Authenticate(w http.ResponseWriter, r *http.Request) {
	defer metric.RecordTimeMetric(time.Now(), "Authenticate")

	id := r.Context().Value(entity.UserIdKey).(int64)

	desc, public, err := a.authLogic.Authenticate(r.Context(), id)
	if err != nil {
		a.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		metric.RecordHitMetric(desc.Status, r.URL.Path)
		return
	}

	res, err := json.Marshal(public)
	if err != nil {
		desc = entity.Description{
			Function: "Authenticate",
			Action:   "Marshal",
			Status:   http.StatusInternalServerError}
		a.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	metric.RecordHitMetric(http.StatusOK, r.URL.Path)
	if _, err := w.Write(res); err != nil {
		a.logger.Error(entity.Description{Function: "Authenticate", Action: "Write"}, err)
	}
}
