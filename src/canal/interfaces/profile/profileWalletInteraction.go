package profile

import (
	"encoding/json"
	"net/http"
	profile "server/src/profile/profile/gen"

	"github.com/sirupsen/logrus"
)

func (p *Profile) GetWallets(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int64)

	wallets, err := p.userApp.GetWallets(r.Context(), &profile.UserID{Id: id})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetWallets",
			"action":   "GetWallets",
			"userID":   id,
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(wallets)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetWallets",
			"action":   "Marshal",
			"userID":   id,
			"wallets":  wallets,
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetWallets",
			"action":   "Write",
			"userID":   id,
			"wallet":   wallets,
		}).Error(err)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	wallet.Id = r.Context().Value("id").(int64)

	if _, err = p.userApp.CreateWallet(r.Context(), &wallet); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "SetWallet",
			"action":   "CreateWallet",
			"wallet":   &wallet,
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
