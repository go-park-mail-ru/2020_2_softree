package profile

import (
	"encoding/json"
	"net/http"
	"server/canal/pkg/domain/entity"
	profile "server/profile/pkg/profile/gen"

	"github.com/sirupsen/logrus"
)

func (p *Profile) GetWallets(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(entity.UserIdKey).(int64)

	wallets, err := p.profile.GetWallets(r.Context(), &profile.UserID{Id: id})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetWallets",
			"action":   "GetWallets",
			"userID":   id,
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(wallets.Wallets)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetWallets",
			"action":   "Marshal",
			"userID":   id,
			"wallets":  wallets,
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	p.recordHitMetric(http.StatusOK)
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetWallets",
			"action":   "Write",
			"userID":   id,
			"wallet":   wallets,
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *Profile) SetWallet(w http.ResponseWriter, r *http.Request) {
	var wallet profile.ConcreteWallet
	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "SetWallet",
			"action":   "Decode",
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	wallet.Id = r.Context().Value(entity.UserIdKey).(int64)

	if _, err = p.profile.CreateWallet(r.Context(), &wallet); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "SetWallet",
			"action":   "CreateWallet",
			"wallet":   &wallet,
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p.recordHitMetric(http.StatusOK)
	w.WriteHeader(http.StatusOK)
}
