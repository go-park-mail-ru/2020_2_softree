package profile

import (
	"encoding/json"
	"net/http"
)

func (p *Profile) GetWallet(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(uint64)

	wallet, err := p.userApp.GetWallet(id)
	if err != nil {
		p.log.Info("user id: ", id, ", func: GetWallet, with error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(wallet)
	if err != nil {
		p.log.Info("wallet: ", wallet, ", func: GetWallet, with error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		p.log.Info("func: GetWallet, with error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *Profile) SetWallet(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(uint64)

	wallet, err := p.userApp.GetWallet(id)
	if err != nil {
		p.log.Info("user id: ", id, ", func: GetWallet, with error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(wallet)
	if err != nil {
		p.log.Info("wallet: ", wallet, ", func: GetWallet, with error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		p.log.Info("func: GetWallet, with error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
