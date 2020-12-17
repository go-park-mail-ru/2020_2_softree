package profile

import (
	json "github.com/mailru/easyjson"
	"net/http"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/infrastructure/metric"
	"time"
)

func (p *Profile) UpdateUserAvatar(w http.ResponseWriter, r *http.Request) {
	defer metric.RecordTimeMetric(time.Now(), "UpdateUserAvatar")

	user, desc, err := entity.GetUserFromBody(r.Body)
	if err != nil {
		desc.Function = "UpdateUserAvatar"
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}
	user.Id = r.Context().Value(entity.UserIdKey).(int64)

	desc, public, err := p.profileLogic.UpdateAvatar(r.Context(), user)
	if err != nil {
		p.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		metric.RecordHitMetric(desc.Status, r.URL.Path)
		return
	}

	res, err := json.Marshal(public)
	if err != nil {
		desc = entity.Description{Function: "UpdateUserAvatar", Action: "Marshal", Status: http.StatusInternalServerError}
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	metric.RecordHitMetric(http.StatusOK, r.URL.Path)
	if _, err := w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "UpdateUserAvatar", Action: "Write"}, err)
	}
}

func (p *Profile) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	defer metric.RecordTimeMetric(time.Now(), "UpdateUserPassword")

	user, desc, err := entity.GetUserFromBody(r.Body)
	if err != nil {
		desc.Function = "UpdateUserPassword"
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}
	user.Id = r.Context().Value(entity.UserIdKey).(int64)

	desc, public, err := p.profileLogic.UpdatePassword(r.Context(), user)
	if err != nil {
		p.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		metric.RecordHitMetric(desc.Status, r.URL.Path)
		return
	}
	if desc.ErrorJSON.NotEmpty {
		metric.RecordHitMetric(p.createServerError(&desc.ErrorJSON, w), r.URL.Path)
		return
	}

	res, err := json.Marshal(public)
	if err != nil {
		desc = entity.Description{Function: "UpdateUserPassword", Action: "Marshal", Status: http.StatusInternalServerError}
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	metric.RecordHitMetric(http.StatusOK, r.URL.Path)
	if _, err := w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "UpdateUserPassword", Action: "Write"}, err)
	}
}

func (p *Profile) GetUser(w http.ResponseWriter, r *http.Request) {
	defer metric.RecordTimeMetric(time.Now(), "GetUser")

	id := r.Context().Value(entity.UserIdKey).(int64)

	desc, public, err := p.profileLogic.ReceiveUser(r.Context(), id)
	if err != nil {
		p.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		metric.RecordHitMetric(desc.Status, r.URL.Path)
		return
	}

	res, err := json.Marshal(public)
	if err != nil {
		desc = entity.Description{Function: "GetUser", Action: "Marshal", Status: http.StatusInternalServerError}
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	metric.RecordHitMetric(http.StatusOK, r.URL.Path)
	if _, err := w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "GetUser", Action: "Write"}, err)
	}
}

func (p *Profile) GetUserWatchlist(w http.ResponseWriter, r *http.Request) {
	defer metric.RecordTimeMetric(time.Now(), "GetUserWatchlist")

	id := r.Context().Value(entity.UserIdKey).(int64)

	desc, currencies, err := p.profileLogic.ReceiveWatchlist(r.Context(), id)
	if err != nil {
		p.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		metric.RecordHitMetric(desc.Status, r.URL.Path)
		return
	}

	res, err := json.Marshal(currencies)
	if err != nil {
		desc = entity.Description{Function: "GetUserWatchlist", Action: "Marshal", Status: http.StatusInternalServerError}
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	metric.RecordHitMetric(http.StatusOK, r.URL.Path)
	if _, err := w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "GetUserWatchlist", Action: "Write"}, err)
	}
}
