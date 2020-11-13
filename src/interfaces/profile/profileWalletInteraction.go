package profile

import (
	"encoding/json"
	"net/http"
	"server/src/domain/entity"
)

func (p *Profile) GetWallets(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(uint64)

	wallet, err := p.userApp.GetWallets(id)
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

	var wallet entity.Wallet
	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		p.log.Info("func: SetTransactions, with error while decode json: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err = p.userApp.CreateWallet(id, wallet.Title); err != nil {
		p.log.Info("user id: ", id, ", func: GetWallet, with error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
