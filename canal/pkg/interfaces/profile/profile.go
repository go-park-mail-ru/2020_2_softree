package profile

import (
	json "github.com/mailru/easyjson"
	"net/http"
	"server/canal/pkg/domain/entity"
)

func (p *Profile) UpdateUserAvatar(w http.ResponseWriter, r *http.Request) {
	user, desc, err := entity.GetUserFromBody(r.Body)
	if err != nil {
		desc.Function = "UpdateUserAvatar"
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}
	user.Id = r.Context().Value(entity.UserIdKey).(int64)

	desc, public, err := p.profileLogic.UpdateAvatar(r.Context(), user)
	if err != nil {
		p.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		p.recordHitMetric(desc.Status)
		return
	}

	res, err := json.Marshal(public)
	if err != nil {
		desc = entity.Description{Function: "UpdateUserAvatar", Action: "Marshal", Status: http.StatusInternalServerError}
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	p.recordHitMetric(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "UpdateUserAvatar", Action: "Write"}, err)
	}
}

func (p *Profile) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	user, desc, err := entity.GetUserFromBody(r.Body)
	if err != nil {
		desc.Function = "UpdateUserPassword"
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}
	user.Id = r.Context().Value(entity.UserIdKey).(int64)

	desc, public, err := p.profileLogic.UpdatePassword(r.Context(), user)
	if err != nil {
		p.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		p.recordHitMetric(desc.Status)
		return
	}
	if desc.ErrorJSON.NotEmpty {
		p.createServerError(&desc.ErrorJSON, w)
		return
	}

	res, err := json.Marshal(public)
	if err != nil {
		desc = entity.Description{Function: "UpdateUserPassword", Action: "Marshal", Status: http.StatusInternalServerError}
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	p.recordHitMetric(http.StatusOK)

	if _, err := w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "UpdateUserPassword", Action: "Write"}, err)
	}
}

func (p *Profile) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(entity.UserIdKey).(int64)

	desc, public, err := p.profileLogic.ReceiveUser(r.Context(), id)
	if err != nil {
		p.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		p.recordHitMetric(desc.Status)
		return
	}

	res, err := json.Marshal(public)
	if err != nil {
		desc = entity.Description{Function: "GetUser", Action: "Marshal", Status: http.StatusInternalServerError}
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	p.recordHitMetric(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "GetUser", Action: "Write"}, err)
	}
}

func (p *Profile) GetUserWatchlist(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(entity.UserIdKey).(int64)

	desc, currencies, err := p.profileLogic.ReceiveWatchlist(r.Context(), id)
	if err != nil {
		p.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		p.recordHitMetric(desc.Status)
		return
	}

	res, err := json.Marshal(currencies)
	if err != nil {
		desc = entity.Description{Function: "GetUserWatchlist", Action: "Marshal", Status: http.StatusInternalServerError}
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	p.recordHitMetric(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "GetUserWatchlist", Action: "Write"}, err)
	}
}
