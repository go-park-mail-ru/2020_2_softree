package profile

import (
	json "github.com/mailru/easyjson"
	"net/http"
	"server/canal/pkg/domain/entity"
)

func (p *Profile) UpdateUserAvatar(w http.ResponseWriter, r *http.Request) {
	user, desc := entity.GetUserFromBody(r.Body)
	if desc.Err != nil {
		desc.Function = "UpdateUserAvatar"
		p.logger.Error(desc)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}
	user.Id = r.Context().Value(entity.UserIdKey).(int64)

	desc, public := p.profileLogic.UpdateAvatar(r.Context(), user)
	if desc.Err != nil {
		p.logger.Error(desc)
		w.WriteHeader(desc.Status)

		p.recordHitMetric(desc.Status)
		return
	}

	res, err := json.Marshal(public)
	if err != nil {
		code := http.StatusInternalServerError
		desc := entity.Description{Function: "UpdateUserAvatar", Action: "Marshal", Err: err, Status: code}
		p.logger.Error(desc)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	p.recordHitMetric(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "UpdateUserAvatar", Action: "Write", Err: err})
	}
}

func (p *Profile) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	user, desc := entity.GetUserFromBody(r.Body)
	if desc.Err != nil {
		desc.Function = "UpdateUserPassword"
		p.logger.Error(desc)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}
	user.Id = r.Context().Value(entity.UserIdKey).(int64)

	desc, public := p.profileLogic.UpdatePassword(r.Context(), user)
	if desc.Err != nil {
		p.logger.Error(desc)
		w.WriteHeader(desc.Status)

		p.recordHitMetric(desc.Status)
		return
	}
	if desc.ErrorJSON.NotEmpty {
		p.createServerError(&desc.ErrorJSON, w)
	}

	res, err := json.Marshal(public)
	if err != nil {
		code := http.StatusInternalServerError
		desc := entity.Description{Function: "UpdateUserPassword", Action: "Marshal", Err: err, Status: code}
		p.logger.Error(desc)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	p.recordHitMetric(http.StatusOK)

	if _, err := w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "UpdateUserPassword", Action: "Write", Err: err})
	}
}

func (p *Profile) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(entity.UserIdKey).(int64)

	desc, public := p.profileLogic.ReceiveUser(r.Context(), id)
	if desc.Err != nil {
		p.logger.Error(desc)
		w.WriteHeader(desc.Status)

		p.recordHitMetric(desc.Status)
		return
	}

	res, err := json.Marshal(public)
	if err != nil {
		code := http.StatusInternalServerError
		desc := entity.Description{Function: "GetUser", Action: "Marshal", Err: err, Status: code}
		p.logger.Error(desc)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	p.recordHitMetric(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "GetUser", Action: "Write", Err: err})
	}
}

func (p *Profile) GetUserWatchlist(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(entity.UserIdKey).(int64)

	desc, currencies := p.profileLogic.ReceiveWatchlist(r.Context(), id)
	if desc.Err != nil {
		p.logger.Error(desc)
		w.WriteHeader(desc.Status)

		p.recordHitMetric(desc.Status)
		return
	}

	res, err := json.Marshal(currencies)
	if err != nil {
		code := http.StatusInternalServerError
		desc := entity.Description{Function: "GetUserWatchlist", Action: "Marshal", Err: err, Status: code}
		p.logger.Error(desc)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	p.recordHitMetric(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "GetUser", Action: "Write", Err: err})
	}
}
