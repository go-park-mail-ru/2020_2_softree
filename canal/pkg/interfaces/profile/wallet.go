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
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	p.recordHitMetric(http.StatusOK)

	if _, err = w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetWallets",
			"action":   "Write",
			"userID":   id,
			"wallet":   wallets,
		}).Error(err)
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
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	p.recordHitMetric(http.StatusOK)
}
